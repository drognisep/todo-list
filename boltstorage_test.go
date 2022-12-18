package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
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
