package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"sort"
)

const (
	boltDBFileName = ".tasklist.db"
	exportPattern  = "taskstore_*.snapshot"
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
		return ZeroTask, err
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
		if errors.Is(err, bolthold.ErrNotFound) {
			return ZeroTask, fmt.Errorf("%w: %v", ErrIDNotFound, err)
		}
		return ZeroTask, err
	}
	return task, nil
}

func (b *boltStorage) Delete(id uint64) error {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.store.TxDelete(tx, id, new(Task))
	})
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}

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

	tasks, err := b.Get()
	if err != nil {
		return
	}
	err = writer.Encode(exportModel{Tasks: tasks})
	if err != nil {
		return
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

	reader := json.NewDecoder(in)
	var model exportModel
	if err := reader.Decode(&model); err != nil {
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

	for _, imported := range model.Tasks {
		switch {
		case imported.ID == 0:
			return rollback(fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/unset ID field"))
		case len(imported.Name) == 0:
			return rollback(fmt.Errorf("%w: %v", ErrUnmappedReqdImportField, "missing/empty name field"))
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
