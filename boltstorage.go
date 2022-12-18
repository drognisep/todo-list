package main

import (
	"errors"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"sort"
)

const (
	boltDBFileName = ".tasklist.db"
)

var (
	ErrIDNotFound  = errors.New("specified ID not found")
	ErrAmbiguousID = errors.New("ambiguous ID detected")
	zeroTask       = Task{}
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
		return zeroTask, nil
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
		return zeroTask, nil
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
