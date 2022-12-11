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

	newTask.ID = created.ID
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

	err = tstore.Delete(created.ID)
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

	updated, err := tstore.Update(created.ID, created)
	assert.NoError(t, err)
	assert.Equal(t, created, updated)

	tasks, err = tstore.Get()
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
	assert.Equal(t, created.Name, tasks[0].Name)
}

func TestTaskStorage_Get_WithID(t *testing.T) {
	tstore := newTaskStorage()

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
