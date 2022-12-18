package data

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	boltDBFileName = ".tasklist.db"
)

var _ TaskStorage = (*boltStorage)(nil)

type boltStorage struct {
	store *bolthold.Store
}

func newBoltStorage() (*boltStorage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return newBoltStorageAt(filepath.Join(home, boltDBFileName))
}

func newBoltStorageAt(file string) (*boltStorage, error) {
	store, err := bolthold.Open(file, 0700, &bolthold.Options{})
	if err != nil {
		return nil, err
	}
	return &boltStorage{store: store}, nil
}

var _ sort.Interface = (*taskSorter)(nil)

type taskSorter struct {
	tasks []Task
}

func newTaskSorter(tasks []Task) *taskSorter {
	return &taskSorter{tasks: tasks}
}

func (t *taskSorter) Len() int {
	return len(t.tasks)
}

func (t *taskSorter) Less(i, j int) bool {
	a, b := t.tasks[i], t.tasks[j]

	switch {
	case !a.Done && b.Done:
		return true
	case a.Done && !b.Done:
		return false
	case a.Favorite && !b.Favorite:
		return true
	case !a.Favorite && b.Favorite:
		return false
	case a.Priority > b.Priority:
		return true
	case a.Priority < b.Priority:
		return false
	case a.ID < b.ID:
		return true
	default:
		return false
	}
}

func (t *taskSorter) Swap(i, j int) {
	t.tasks[i], t.tasks[j] = t.tasks[j], t.tasks[i]
}

func (b *boltStorage) Get(filter ...TaskFilter) ([]Task, error) {
	query := new(bolthold.Query)

	for _, f := range filter {
		f(query)
	}

	var found []Task
	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		if err := b.store.TxFind(tx, &found, query); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(newTaskSorter(found))
	return found, nil
}

func (b *boltStorage) Count() (int, error) {
	var count int
	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		var err error
		count, err = b.store.TxCount(tx, new(Task), new(bolthold.Query))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *boltStorage) Create(task Task) (Task, error) {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.store.TxInsert(tx, bolthold.NextSequence(), &task); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ZeroTask, nil
	}
	return task, nil
}

func (b *boltStorage) Update(id uint64, task Task) (Task, error) {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.store.TxUpdate(tx, id, &task); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return ZeroTask, nil
	}
	return task, nil
}

func (b *boltStorage) Delete(id uint64) error {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.store.TxDelete(tx, id, new(Task))
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	exportColumns = []string{"ID", "NAME", "DESCRIPTION", "DONE", "PRIORITY", "FAVORITE"}
)

func (b *boltStorage) Export(dir string) (outName string, err error) {
	out, err := os.CreateTemp(dir, "taskstore_*.snapshot")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = out.Close()
	}()

	fi, err := out.Stat()
	if err != nil {
		return "", err
	}
	outName = filepath.Join(dir, fi.Name())

	writer := csv.NewWriter(out)
	writer.UseCRLF = false
	defer func() {
		writer.Flush()
	}()
	err = writer.Write(exportColumns)
	if err != nil {
		return
	}

	tasks, err := b.Get()
	if err != nil {
		return
	}

	buf := [6]string{}
	for _, t := range tasks {
		buf[0] = strconv.FormatUint(t.ID, 10)
		buf[1] = t.Name
		buf[2] = t.Description
		buf[3] = strconv.FormatBool(t.Done)
		buf[4] = strconv.FormatInt(int64(t.Priority), 10)
		buf[5] = strconv.FormatBool(t.Favorite)
		if err := writer.Write(buf[:]); err != nil {
			return outName, err
		}
	}

	return outName, nil
}

func (b *boltStorage) Import(file string, merge MergeStrategy) error {
	in, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		_ = in.Close()
	}()

	reader := csv.NewReader(in)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	tx, err := b.store.Bolt().Begin(true)
	if err != nil {
		return err
	}
	fn := b.mapInput(header)

	rollback := func(err error) error {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: %v", err, rerr)
		}
		return err
	}

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return rollback(err)
		}
		var imported Task
		imported, err = fn(imported, row)
		if err != nil {
			return rollback(err)
		}

		var dupe Task
		err = b.store.TxFindOne(tx, &dupe, bolthold.Where(bolthold.Key).Eq(imported.ID))
		if err != nil {
			if err != bolthold.ErrNotFound {
				return rollback(err)
			}

			// No dupe found, import it directly.
			if err := b.store.TxInsert(tx, imported.ID, imported); err != nil {
				return rollback(err)
			}
			continue
		}

		// A duplicate was found, employ the MergeStrategy.
		switch merge {
		case MergeOverwrite:
			dupe.Name = imported.Name
			dupe.Description = imported.Description
			dupe.Done = imported.Done
			dupe.Priority = imported.Priority
			dupe.Favorite = imported.Favorite

			err := b.store.TxUpdate(tx, dupe.ID, dupe)
			if err != nil {
				return rollback(err)
			}
			continue
		case MergeKeepInternal:
			continue
		case MergeError:
			return rollback(errors.New("conflict discovered, rolling back import"))
		default:
			return rollback(fmt.Errorf("unrecognized merge strategy '%s'", merge))
		}
	}

	if err := tx.Commit(); err != nil {
		return rollback(err)
	}
	return nil
}

type mapFunc func(task Task, row []string) (Task, error)

func (b *boltStorage) mapInput(header []string) mapFunc {
	inputIdx := [6]int{-1, -1, -1, -1, -1, -1}

	for i, base := range exportColumns {
		for j, in := range header {
			if base == in {
				inputIdx[i] = j
				break
			}
		}
	}

	return func(task Task, row []string) (Task, error) {
		for i, idx := range inputIdx {
			switch i {
			case 0:
				if idx == -1 {
					return ZeroTask, fmt.Errorf("%w: %s", ErrUnmappedReqdImportField, "missing ID field")
				}
				id, err := strconv.ParseUint(row[idx], 10, 64)
				if err != nil {
					return ZeroTask, fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, err)
				}
				task.ID = id
			case 1:
				if idx == -1 {
					return ZeroTask, fmt.Errorf("%w: %s", ErrUnmappedReqdImportField, "missing name field")
				}
				if len(row[idx]) == 0 {
					return ZeroTask, fmt.Errorf("%w: %s", ErrUnmappedReqdImportField, "empty name field")
				}
				task.Name = row[idx]
			case 2:
				if idx == -1 {
					task.Description = ""
					continue
				}
				task.Description = row[idx]
			case 3:
				if idx == -1 {
					task.Done = false
					continue
				}
				b, err := strconv.ParseBool(row[idx])
				if err != nil {
					task.Done = false
					continue
				}
				task.Done = b
			case 4:
				if idx == -1 {
					task.Priority = 0
					continue
				}
				prio, err := strconv.ParseInt(row[idx], 10, 32)
				if err != nil {
					task.Priority = 0
					continue
				}
				task.Priority = int(prio)
			case 5:
				if idx == -1 {
					task.Favorite = false
					continue
				}
				tf, err := strconv.ParseBool(row[idx])
				if err != nil {
					task.Favorite = false
					continue
				}
				task.Favorite = tf
			}
		}
		return task, nil
	}
}
