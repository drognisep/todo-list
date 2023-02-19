package data

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	created, err := tstore.CreateTask(newTask)
	assert.NoError(t, err)
	assert.NotEqual(t, created, newTask)

	newTask.ID = created.ID
	assert.Equal(t, created, newTask)

	tasks, err := tstore.GetTasks()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created, tasks[0])
}

func TestBoltStorage_Delete(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	tasks, err := tstore.GetTasks()
	assert.NoError(t, err)
	assert.Empty(t, tasks)

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.CreateTask(newTask)
	assert.NoError(t, err)

	tasks, err = tstore.GetTasks()
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks)

	err = tstore.DeleteTask(created.ID)
	assert.NoError(t, err)

	tasks, err = tstore.GetTasks()
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestBoltStorage_Update(t *testing.T) {
	tstore, cleanup := _newBoltStore(t)
	defer cleanup()

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.CreateTask(newTask)
	assert.NoError(t, err)
	assert.NotEqual(t, created, newTask)

	created.Name = "New name"

	tasks, err := tstore.GetTasks()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.NotEqual(t, created.Name, tasks[0].Name)

	updated, err := tstore.UpdateTask(created.ID, created)
	assert.NoError(t, err)
	assert.Equal(t, created, updated)

	tasks, err = tstore.GetTasks()
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
	_, err = tstore.CreateTask(newTask)
	assert.NoError(t, err)
	_, err = tstore.CreateTask(newTask)
	assert.NoError(t, err)
	_, err = tstore.CreateTask(newTask)
	assert.NoError(t, err)

	count, err := tstore.CountTasks()
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	tasks, err := tstore.GetTasks(WithID(1))
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, uint64(1), tasks[0].ID)
}

func TestBoltStorage_StartTimeEntry(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	var err error
	task := Task{Name: "First Task"}
	task, err = store.CreateTask(task)
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
	task, err = store.CreateTask(task)
	require.NoError(t, err)
	require.Equal(t, uint64(1), task.ID)

	entry, err := store.StartTimeEntry(task.ID)
	require.NoError(t, err)

	assert.NoError(t, store.StopTimeEntry(entry.ID))
	entries, err := store.GetTimeEntries()
	require.Len(t, entries, 1)
	require.NotNil(t, entries[0].End)
	assert.LessOrEqual(t, *entries[0].End, time.Now())
}

func TestBoltStorage_StartAfterStop(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	var err error
	task := Task{Name: "First Task"}
	task, err = store.CreateTask(task)
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

func TestBoltStorage_GetRunningTimeEntry(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{
		Name: "A task",
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), created.ID)

	entry, err := store.StartTimeEntry(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), entry.ID)
	assert.Equal(t, created.ID, entry.TaskID)

	running, err := store.GetRunningTimeEntry()
	assert.NoError(t, err)
	require.NotNil(t, running)
	assert.Equal(t, entry.ID, running.ID)
}

func TestBoltStorage_GetTimeEntries_After(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{
		Name: "Task",
	})
	require.NoError(t, err)

	entry, err := store.StartTimeEntry(created.ID)
	require.NoError(t, err)
	require.NoError(t, store.StopTimeEntry(entry.ID))

	entries, err := store.GetTimeEntries(SinceWeekday(time.Sunday))
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, entry.ID, entries[0].ID)
}

func TestBoltStorage_GetTimeEntries_EntriesToday(t *testing.T) {
	store, cleanup := _newBoltStore(t)
	defer cleanup()

	created, err := store.CreateTask(Task{
		Name: "Task",
	})
	require.NoError(t, err)

	entry, err := store.StartTimeEntry(created.ID)
	require.NoError(t, err)
	require.NoError(t, store.StopTimeEntry(entry.ID))

	entries, err := store.GetTimeEntries(SinceWeekday(time.Sunday))
	assert.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, entry.ID, entries[0].ID)
}

func _newBoltStore(t *testing.T) (*boltStorage, func()) {
	testName := strings.ReplaceAll(strings.ReplaceAll(t.Name(), "/", "__"), "\\", "__")
	temp, err := os.MkdirTemp("", testName+"_*")
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
		} else {
			t.Logf("Cleaned up test DB '%s'", temp)
		}
	}
}
