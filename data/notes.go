package data

import (
	"errors"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
	"time"
)

var (
	ErrEmptyNote = errors.New("note text cannot be empty")
)

type Note struct {
	ID      uint64    `json:"id" boltholdKey:"ID"`
	TaskID  *uint64   `json:"taskID" boltholdIndex:"Note_Task"`
	Text    string    `json:"text"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (b *boltStorage) AddNote(id uint64, note Note) (Note, error) {
	trimmedText := strings.TrimSpace(note.Text)
	if len(trimmedText) == 0 {
		return Note{}, ErrEmptyNote
	}
	note.Text = trimmedText
	note.ID = 0
	now := time.Now()
	note.Created = now
	note.Updated = now
	note.TaskID = &id
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.store.TxInsert(tx, bolthold.NextSequence(), &note); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

var _ sort.Interface = (noteSorter)(nil)

type noteSorter []Note

func (n noteSorter) Len() int {
	return len(n)
}

func (n noteSorter) Less(i, j int) bool {
	switch {
	case n[i].Created.Before(n[j].Created):
		return true
	}
	return false
}

func (n noteSorter) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (b *boltStorage) GetTaskNotes(id uint64) ([]Note, error) {
	var notes []Note
	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		err := b.store.TxFind(tx, &notes, bolthold.Where("TaskID").MatchFunc(func(note *Note) (bool, error) {
			if note.TaskID != nil && *note.TaskID == id {
				return true, nil
			}
			return false, nil
		}).Index("Note_Task"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Sort(noteSorter(notes))
	return notes, nil
}

func (b *boltStorage) GetTaskNoteCount(id uint64) (int, error) {
	var count int
	err := b.store.Bolt().View(func(tx *bbolt.Tx) error {
		var err error
		count, err = b.store.TxCount(tx, Note{}, bolthold.Where("TaskID").MatchFunc(func(note *Note) (bool, error) {
			if note.TaskID != nil && *note.TaskID == id {
				return true, nil
			}
			return false, nil
		}).Index("Note_Task"))
		return err
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (b *boltStorage) UpdateNote(id uint64, note Note) (Note, error) {
	trimmedText := strings.TrimSpace(note.Text)
	if len(trimmedText) == 0 {
		return Note{}, ErrEmptyNote
	}
	note.Text = trimmedText
	var (
		existing Note
	)

	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		if err := b.store.TxFindOne(tx, &existing, bolthold.Where(bolthold.Key).Eq(id)); err != nil {
			return translateNotFound(err, id)
		}
		existing.Text = note.Text
		existing.Updated = time.Now()
		if err := b.store.TxUpdate(tx, id, &existing); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Note{}, err
	}
	return existing, nil
}

func (b *boltStorage) DeleteNote(id uint64) error {
	err := b.store.Bolt().Update(func(tx *bbolt.Tx) error {
		return b.store.TxDelete(tx, id, Note{})
	})
	if err != nil {
		return err
	}
	return nil
}
