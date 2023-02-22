package data

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAddNoteToTask(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{Name: "Task with notes"})
	require.NoError(t, err)

	var newNote Note
	note := Note{
		Text: "A note about a task with notes",
	}
	newNote, err = store.AddNote(created.ID, note)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), newNote.ID)
	assert.NotNil(t, newNote.TaskID)
	assert.Equal(t, created.ID, *newNote.TaskID)
	assert.Equal(t, "A note about a task with notes", newNote.Text)
	assert.True(t, time.Now().After(newNote.Created), "Created notes should have a create datetime before now")
	assert.True(t, time.Now().After(newNote.Updated), "Created notes should have a create datetime before now")
	assert.Equal(t, newNote.Created, newNote.Updated)
}

func TestAddNoteToTask_Neg(t *testing.T) {
	tests := map[string]struct {
		note        Note
		expectedErr error
	}{
		"Missing text": {
			note:        Note{Text: "\r\n   \t\n"},
			expectedErr: ErrEmptyNote,
		},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			store, cleanup := _newBoltStore(t)
			defer cleanup()

			newTask, err := store.CreateTask(Task{Name: "test task"})
			require.NoError(t, err)

			newNote, err := store.AddNote(newTask.ID, tc.note)
			assert.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, uint64(0), newNote.ID)
		})
	}
}

func TestGetTaskNotes(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{Name: "Task with notes"})
	require.NoError(t, err)

	note := Note{Text: "testing"}
	newNote, err := store.AddNote(created.ID, note)
	require.NoError(t, err)

	var notes []Note
	notes, err = store.GetTaskNotes(created.ID)
	assert.NoError(t, err)
	assert.Len(t, notes, 1)
	assert.Equal(t, newNote.ID, notes[0].ID)
}

func TestGetTaskNoteCount(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{Name: "Task with notes"})
	require.NoError(t, err)

	note := Note{Text: "testing"}
	newNote, err := store.AddNote(created.ID, note)
	require.NoError(t, err)

	var count int
	count, err = store.GetTaskNoteCount(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	err = store.DeleteNote(newNote.ID)
	assert.NoError(t, err)

	count, err = store.GetTaskNoteCount(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestUpdateTaskNote(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	task, err := store.CreateTask(Task{Name: "new task"})
	require.NoError(t, err)

	note, err := store.AddNote(task.ID, Note{Text: "new note"})
	require.NoError(t, err)
	require.Equal(t, uint64(1), note.ID)

	var updated Note
	note.Text = "updated text"
	updated, err = store.UpdateNote(note.ID, note)
	assert.NoError(t, err)
	assert.Equal(t, note.ID, updated.ID)
	assert.Equal(t, "updated text", updated.Text)
	assert.NotEqual(t, "new note", updated.Text)
	assert.NotEqual(t, note.Updated, updated.Updated)
	assert.True(t, updated.Created.Before(updated.Updated), "created should be before updated")
}

func TestDeleteNote(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	task, err := store.CreateTask(Task{Name: "new task"})
	require.NoError(t, err)

	note, err := store.AddNote(task.ID, Note{Text: "new note"})
	require.NoError(t, err)
	require.Equal(t, uint64(1), note.ID)

	err = store.DeleteNote(note.ID)
	assert.NoError(t, err)

	notes, err := store.GetTaskNotes(task.ID)
	assert.NoError(t, err)
	assert.Empty(t, notes)
}
