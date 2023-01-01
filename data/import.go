package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"os"
	"reflect"
	"runtime/debug"
	"time"
)

func (b *boltStorage) Export(dir string) (outName string, err error) {
	out, err := os.CreateTemp(dir, exportPattern)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()

	outName = out.Name()

	writer := json.NewEncoder(out)

	tasks, err := b.GetHistoric()
	if err != nil {
		return "", err
	}
	entries, err := b.GetTimeEntries()
	if err != nil {
		return "", err
	}
	err = writer.Encode(exportModel{Tasks: tasks, TimeEntries: entries})
	if err != nil {
		return
	}

	return outName, nil
}

func (b *boltStorage) Import(file string, merge MergeStrategy) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	var model exportModel
	if err := json.NewDecoder(in).Decode(&model); err != nil {
		return err
	}

	tx, err := b.store.Bolt().Begin(true)
	if err != nil {
		return err
	}
	rollback := func(err error) error {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: %v", err, rerr)
		}
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			rerr := fmt.Errorf("%v", r)
			if err != nil {
				err = rollback(fmt.Errorf("%w: %v", err, rerr))
			}
			err = rollback(rerr)
		}
	}()

	taskSet := &importSet[Task]{
		records: model.Tasks,
		getID:   func(t Task) uint64 { return t.ID },
		validateImport: func(t Task) error {
			switch {
			case len(t.Name) == 0:
				return fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/empty name field")
			}
			return nil
		},
	}
	timeEntrySet := &importSet[TimeEntry]{
		records: model.TimeEntries,
		getID:   func(e TimeEntry) uint64 { return e.ID },
		validateImport: func(e TimeEntry) error {
			switch {
			case e.TaskID == 0:
				return fmt.Errorf("%w: invalid Task ID of 0", ErrUnmappedReqdImportField)
			case e.Start == time.Time{}:
				return fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/empty start field")
			}

			count, err := b.store.TxCount(tx, Task{}, bolthold.Where(bolthold.Key).Eq(e.TaskID))
			switch {
			case err != nil:
				return err
			case count > 1:
				return ErrAmbiguousID
			case count == 0:
				return fmt.Errorf("%w: referenced Task ID '%d' does not exist", ErrIDNotFound, e.TaskID)
			default:
				return nil
			}
		},
	}
	if err := importRecords(b, tx, taskSet, merge); err != nil {
		return rollback(err)
	}
	if err := importRecords(b, tx, timeEntrySet, merge); err != nil {
		return rollback(err)
	}

	if err := tx.Commit(); err != nil {
		return rollback(err)
	}
	return nil
}

type importSet[T comparable] struct {
	records        []T
	getID          func(T) uint64
	validateImport func(T) error
}

func importRecords[T comparable](b *boltStorage, tx *bbolt.Tx, dataSet *importSet[T], merge MergeStrategy) error {
	if len(dataSet.records) == 0 {
		return nil
	}

	for _, imported := range dataSet.records {
		if dataSet.getID(imported) == 0 {
			return fmt.Errorf("%w: missing/unset ID field", ErrUnmappedReqdImportField)
		}
		if err := dataSet.validateImport(imported); err != nil {
			return err
		}

		var (
			dupe T
			err  error
		)
		err = b.store.TxFindOne(tx, &dupe, bolthold.Where(bolthold.Key).Eq(dataSet.getID(imported)))
		if err != nil {
			if err != bolthold.ErrNotFound {
				return err
			}

			// No dupe found, import it directly.
			if err := b.store.TxInsert(tx, dataSet.getID(imported), imported); err != nil {
				return err
			}
			continue
		}

		// Check if the dupe represents a different state.
		if dupe == imported {
			continue
		}

		// A duplicate was found, employ the MergeStrategy.
		switch merge {
		case MergeOverwrite:
			err := b.store.TxUpdate(tx, dataSet.getID(dupe), imported)
			if err != nil {
				return err
			}
			continue
		case MergeKeepInternal:
			continue
		case MergeError:
			return errors.New("conflict discovered, rolling back import")
		default:
			return fmt.Errorf("unrecognized merge strategy '%s'", merge)
		}
	}

	// Reset key sequence after all imports have been added.
	var (
		record T
		max    uint64
	)
	agg, err := b.store.TxFindAggregate(tx, record, nil)
	if err != nil {
		return err
	}
	for _, a := range agg {
		a.Max("ID", &record)
		_max := dataSet.getID(record)
		if _max > max {
			max = _max
		}
	}
	typeName := reflect.TypeOf(record).Name()
	bucket, err := tx.CreateBucketIfNotExists([]byte(typeName))
	if err != nil {
		return fmt.Errorf("%w: unable to create bucket for type '%s'", err, typeName)
	}
	return bucket.SetSequence(max)
}
