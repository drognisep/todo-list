package data

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

func TestNewBoltStorage(t *testing.T) {
	_, cleanup := _newBoltStore(t)
	cleanup()
}

func TestBoltStorage_Create(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.Create(newTask)
	assert.NoError(t, err)
	assert.NotEqual(t, created, newTask)

	newTask.ID = created.ID
	assert.Equal(t, created, newTask)

	tasks, err := tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created, tasks[0])
}

func TestBoltStorage_Delete(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	tasks, err := tstore.Get()
	assert.NoError(t, err)
	assert.Empty(t, tasks)

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.Create(newTask)
	assert.NoError(t, err)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)

	err = tstore.Delete(created.ID)
	assert.NoError(t, err)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestBoltStorage_Update(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.Create(newTask)
	assert.NoError(t, err)
	assert.NotEqual(t, created, newTask)

	created.Name = "New name"

	tasks, err := tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.NotEqual(t, created.Name, tasks[0].Name)

	updated, err := tstore.Update(created.ID, created)
	assert.NoError(t, err)
	assert.Equal(t, created, updated)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created.Name, tasks[0].Name)
}

func TestBoltStorage_Get_WithID(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	newTask := Task{
		Name: "Some name",
	}

	var err error
	_, err = tstore.Create(newTask)
	assert.NoError(t, err)
	_, err = tstore.Create(newTask)
	assert.NoError(t, err)
	_, err = tstore.Create(newTask)
	assert.NoError(t, err)

	count, err := tstore.Count()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	tasks, err := tstore.Get(WithID(1))
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, uint64(1), tasks[0].ID)
}

func TestBoltStorage_Export(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.Create(newTask)
	assert.NoError(t, err)

	temp, err := os.MkdirTemp("", "export_*")
	assert.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(temp)
	}()

	snapshot, err := tstore.Export(temp)
	require.NoError(t, err)

	assert.NoError(t, tstore.Delete(created.ID))
	count, err := tstore.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Soft deleted task should not show up with count")

	require.NoError(t, tstore.Import(snapshot, MergeOverwrite))
	count, err = tstore.Count()
	assert.NoError(t, err)
	require.Equal(t, 1, count, "Snapshot will overwrite state of the soft-deleted task")

	found, err := tstore.Get(WithID(created.ID))
	assert.NoError(t, err)
	assert.Equal(t, created, found[0])
}

func TestBoltStorage_Import(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	data := `{
"tasks": [{
	"name": "Test Task",
	"id": 1,
	"done": false,
	"description": "With a description",
	"favorite": true,
	"priority": 3
}]
}`

	temp, err := os.CreateTemp("", "import.csv")
	assert.NoError(t, err)
	defer func() {
		_ = temp.Close()
		assert.NoError(t, os.Remove(temp.Name()))
	}()

	_, err = io.Copy(temp, strings.NewReader(data))
	assert.NoError(t, err)

	require.NoError(t, tstore.Import(temp.Name(), MergeError))

	found, err := tstore.Get()
	assert.NoError(t, err)
	assert.Equal(t, Task{
		ID:          1,
		Name:        "Test Task",
		Description: "With a description",
		Done:        false,
		Favorite:    true,
		Priority:    3,
	}, found[0])
}

func TestBoltStorage_Import_ID_not_found(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	data := `{
"tasks": [{
	"name": "Test Task",
	"done": false,
	"description": "With a description",
	"favorite": true,
	"priority": 3
}]
}`

	temp, err := os.CreateTemp("", "import.csv")
	assert.NoError(t, err)
	defer func() {
		_ = temp.Close()
		assert.NoError(t, os.Remove(temp.Name()))
	}()

	_, err = io.Copy(temp, strings.NewReader(data))
	assert.NoError(t, err)

	err = tstore.Import(temp.Name(), MergeError)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrUnmappedReqdImportField)
}

func TestBoltStorage_StartTimeEntry(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	var err error
	task := Task{Name: "First Task"}
	task, err = store.Create(task)
	assert.NoError(t, err)

	entry, err := store.StartTimeEntry(task.ID)
	assert.NoError(t, err)
	assert.Equal(t, task.ID, entry.TaskID)
	assert.Less(t, time.Time{}, entry.Start)
}

func TestBoltStorage_StopTimeEntry(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	var err error
	task := Task{Name: "First Task"}
	task, err = store.Create(task)
	require.NoError(t, err)
	require.Equal(t, uint64(1), task.ID)

	entry, err := store.StartTimeEntry(task.ID)
	require.NoError(t, err)

	assert.NoError(t, store.StopTimeEntry(entry.ID))
	entries, err := store.GetTimeEntries()
	require.Len(t, entries, 1)
	require.NotNil(t, entries[0].End)
	assert.Less(t, *entries[0].End, time.Now())
}

func TestBoltStorage_StartAfterStop(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	var err error
	task := Task{Name: "First Task"}
	task, err = store.Create(task)
	require.NoError(t, err)
	require.Equal(t, uint64(1), task.ID)

	firstEntry, err := store.StartTimeEntry(task.ID)
	require.NoError(t, err)
	time.Sleep(200 * time.Millisecond)

	secondEntry, err := store.StartAfterStop(task.ID, firstEntry.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, firstEntry.ID, secondEntry.ID)

	entries, err := store.GetTimeEntries()
	require.Len(t, entries, 2)
	assert.Equal(t, task.ID, entries[0].TaskID)
	assert.NotNil(t, entries[0].End)
	assert.Equal(t, task.ID, entries[1].TaskID)
	assert.Nil(t, entries[1].End)
}

func _newBoltStore(t *testing.T) (*boltStorage, func()) {
	temp, err := os.MkdirTemp("", t.Name()+"_*")
	require.NoError(t, err, "Failed to create temp dir")

	dbFile := filepath.Join(temp, "testDB.db")
	t.Logf("Test DB file for '%s' is '%s'", t.Name(), dbFile)
	storage, err := newBoltStorageAt(dbFile)
	if err != nil {
		assert.NoError(t, os.RemoveAll(temp))
		t.Fatal("Failed to open DB file", err)
	}
	return storage, func() {
		if r := recover(); r != nil {
			t.Log("Recovered panic", r)
			t.Fail()
			debug.PrintStack()
		}
		assert.NotPanics(t, func() {
			if err := storage.store.Close(); err != nil {
				t.Log("Failed to close DB", err)
				t.Fail()
			}
		})
		if err := os.RemoveAll(temp); err != nil {
			t.Log("Failed to cleanup temp dir", err)
			t.Fail()
		}
	}
}
