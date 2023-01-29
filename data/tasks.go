package data

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"sort"
)

type Task struct {
	ID          uint64 `json:"id" boltholdKey:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done" boltholdIndex:"Done"`
	Priority    int    `json:"priority" boltholdIndex:"Priority"`
	Favorite    bool   `json:"favorite"`
	SoftDeleted bool   `json:"inactivated" boltholdIndex:"SoftDeleted"`
}

type exportModel struct {
	Tasks       []Task      `json:"tasks"`
	TimeEntries []TimeEntry `json:"timeEntries,omitempty"`
	Notes       []Note      `json:"notes,omitempty"`
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

func (b *boltStorage) Get(filters ...TaskFilter) ([]Task, error) {
	query := bolthold.Where("SoftDeleted").Eq(false)
	return b.getTasks(query, filters)
}

func (b *boltStorage) GetHistoric(filters ...TaskFilter) ([]Task, error) {
	query := new(bolthold.Query)
	return b.getTasks(query, filters)
}

func (b *boltStorage) getTasks(query *bolthold.Query, filters []TaskFilter) ([]Task, error) {
	for _, f := range filters {
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
		count, err = b.store.TxCount(tx, new(Task), bolthold.Where("SoftDeleted").Eq(false))
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
		return Task{}, err
	}
	return task, nil
}

func (b *boltStorage) Update(id uint64, task Task) (Task, error) {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.store.TxUpdate(tx, id, &task)
	})
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return Task{}, fmt.Errorf("%w: %v", ErrIDNotFound, err)
		}
		return Task{}, err
	}
	return task, nil
}

func (b *boltStorage) Delete(id uint64) error {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		task := new(Task)
		if err := b.store.TxFindOne(tx, task, bolthold.Where(bolthold.Key).Eq(id)); err != nil {
			return nil
		}
		task.SoftDeleted = true
		return b.store.TxUpdate(tx, id, task)
	})
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}
