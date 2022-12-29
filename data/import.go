package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"os"
	"time"
)

func (b *boltStorage) Import(file string, merge MergeStrategy) error {
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

	taskSet := &importSet[Task]{
		records: model.Tasks,
		getID:   func(t Task) uint64 { return t.ID },
		validateImport: func(t Task) error {
			switch {
			case t.ID == 0:
				return fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/unset ID field")
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
			case e.ID == 0:
				return fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/unset ID field")
			case e.Start == time.Time{}:
				return fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/empty start field")
			}
			return nil
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

	for _, imported := range dataSet.records {
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
	return nil
}
