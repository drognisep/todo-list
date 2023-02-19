package data

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"time"
)

type TimeEntry struct {
	ID     uint64     `json:"id" boltholdKey:"ID"`
	TaskID uint64     `json:"taskID" boltholdIndex:"TaskID"`
	Start  time.Time  `json:"start" boltholdIndex:"Start"`
	End    *time.Time `json:"end"`
	Synced bool       `json:"synced" boltholdIndex:"Synced"`
}

func (e *TimeEntry) Duration() time.Duration {
	if e.End == nil {
		return time.Duration(0)
	}
	return e.End.Sub(e.Start)
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

func entryID(id uint64) TimeEntryFilter {
	return func(query *bolthold.Query) {
		query.And(bolthold.Key).Eq(id)
	}
}

func (b *boltStorage) GetTimeEntry(id uint64) (TimeEntry, error) {
	entries, err := b.GetTimeEntries(entryID(id))
	if err != nil {
		return TimeEntry{}, err
	}
	if len(entries) == 0 {
		return TimeEntry{}, fmt.Errorf("%w: %v", ErrIDNotFound, err)
	}
	return entries[0], nil
}

func (b *boltStorage) GetTimeEntries(filters ...TimeEntryFilter) ([]TimeEntry, error) {
	var entries []TimeEntry
	query := new(bolthold.Query)

	for _, filter := range filters {
		filter(query)
	}

	query.SortBy("Start")

	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		return b.store.TxFind(tx, &entries, query)
	})
	if err != nil {
		return nil, err
	}
	return entries, err
}

func (b *boltStorage) GetRunningTimeEntry() (*TimeEntry, error) {
	var result TimeEntry
	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		return b.store.TxFindOne(tx, &result,
			bolthold.Where("End").
				MatchFunc(func(entry *TimeEntry) (bool, error) {
					return entry.End == nil, nil
				}).
				SortBy("Start").Reverse())
	})
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (b *boltStorage) UpdateTimeEntry(id uint64, newState TimeEntry) (TimeEntry, error) {
	var current TimeEntry
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		err := b.store.TxFindOne(tx, &current, bolthold.Where(bolthold.Key).Eq(id))
		if err != nil {
			return err
		}
		current.Start = newState.Start
		if newState.End != nil {
			current.End = newState.End
		}
		current.Synced = false

		if err = b.store.TxUpdate(tx, id, &current); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TimeEntry{}, err
	}
	return current, nil
}
