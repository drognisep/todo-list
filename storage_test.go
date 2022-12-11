package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaskStorage_Create(t *testing.T) {
	tstore := newTaskStorage()

	newTask := Task{
		Name: "Some name",
	}

	created, err := tstore.Create(newTask)
	assert.NoError(t, err)
	assert.NotEqual(t, created, newTask)

	newTask.TaskID = created.TaskID
	assert.Equal(t, created, newTask)

	tasks, err := tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created, tasks[0])
}

func TestTaskStorage_Delete(t *testing.T) {
	tstore := newTaskStorage()

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

	err = tstore.Delete(created.TaskID)
	assert.NoError(t, err)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

func TestTaskStorage_Update(t *testing.T) {
	tstore := newTaskStorage()

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

	updated, err := tstore.Update(created.TaskID, created)
	assert.NoError(t, err)
	assert.Equal(t, created, updated)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created.Name, tasks[0].Name)
}
