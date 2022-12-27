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
	"time"
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
	store, err := bolthold.Open(file, 0700, &bolthold.Options{
		Encoder: func(value any) ([]byte, error) {
			return json.Marshal(value)
		},
		Decoder: func(data []byte, value any) error {
			return json.Unmarshal(data, value)
		},
	})
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

		// Check if the duple represents a different state.
		if dupe == imported {
			continue
		}

		// A duplicate was found, employ the MergeStrategy.
		switch merge {
		case MergeOverwrite:
			err := b.store.TxUpdate(tx, dupe.ID, imported)
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

func (b *boltStorage) StartTimeEntry(taskID uint64) (TimeEntry, error) {
	var entry TimeEntry
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		now := time.Now()
		var err error
		entry, err = b.txStartTimeEntry(tx, taskID, now)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (b *boltStorage) txStartTimeEntry(tx *bbolt.Tx, taskID uint64, now time.Time) (TimeEntry, error) {
	task := new(Task)
	if err := b.store.TxFindOne(tx, task, bolthold.Where(bolthold.Key).Eq(taskID)); err != nil {
		return TimeEntry{}, err
	}
	entry := TimeEntry{
		TaskID: taskID,
		Start:  now,
	}
	if err := b.store.TxInsert(tx, bolthold.NextSequence(), &entry); err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (b *boltStorage) StopTimeEntry(entryID uint64) error {
	now := time.Now()

	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.txStopTimeEntry(tx, entryID, now)
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *boltStorage) txStopTimeEntry(tx *bbolt.Tx, entryID uint64, now time.Time) error {
	found := new(TimeEntry)
	if err := b.store.TxFindOne(tx, found, bolthold.Where(bolthold.Key).Eq(entryID)); err != nil {
		return err
	}
	found.End = &now
	return b.store.TxUpdate(tx, entryID, *found)
}

func (b *boltStorage) StartAfterStop(startTaskID uint64, stopEntryID uint64) (TimeEntry, error) {
	now := time.Now()
	var entry TimeEntry

	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.txStopTimeEntry(tx, stopEntryID, now); err != nil {
			return err
		}
		_entry, err := b.txStartTimeEntry(tx, startTaskID, now)
		if err != nil {
			return err
		}
		entry = _entry
		return nil
	})
	return entry, err
}

func (b *boltStorage) GetTimeEntries(filters ...TimeEntryFilter) ([]TimeEntry, error) {
	var entries []TimeEntry
	query := new(bolthold.Query)

	for _, filter := range filters {
		filter(query)
	}

	query.SortBy("Start")

	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.store.TxFind(tx, &entries, query)
	})
	if err != nil {
		return nil, err
	}
	return entries, err
}
