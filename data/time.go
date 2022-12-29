package data

import (
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
