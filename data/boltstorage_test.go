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

	snapshot, err := tstore.Export(temp)
	require.NoError(t, err)

	assert.NoError(t, tstore.Delete(created.ID))
	count, err := tstore.Count()
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	require.NoError(t, tstore.Import(snapshot, MergeError))
	count, err = tstore.Count()
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

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

func _newBoltStore(t *testing.T) (*boltStorage, func()) {
	temp, err := os.MkdirTemp("", t.Name()+"_*")
	if err != nil {
		t.Fatal("Failed to create temp dir", err)
	}

	dbFile := filepath.Join(temp, "testDB.db")
	t.Logf("Test DB file for '%s' is '%s'", t.Name(), dbFile)
	storage, err := newBoltStorageAt(dbFile)
	if err != nil {
		t.Fatal("Failed to open DB file", err)
	}
	return storage, func() {
		if r := recover(); r != nil {
			t.Log("Recovered panic", r)
			t.Fail()
			debug.PrintStack()
		}
		if err := storage.store.Close(); err != nil {
			t.Log("Failed to close DB", err)
			t.Fail()
		}
		if err := os.RemoveAll(temp); err != nil {
			t.Log("Failed to cleanup temp dir", err)
			t.Fail()
		}
	}
}
