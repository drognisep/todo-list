package main

import (
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"todo-list/data"
)

type TaskController struct {
	store data.TaskStorage
}

func NewTaskController() (*TaskController, error) {
	store, err := data.NewStorage()
	if err != nil {
		return nil, err
	}
	return &TaskController{
		store: store,
	}, nil
}

func (c *TaskController) GetTaskByID(id uint64) (data.Task, error) {
	tasks, err := c.store.Get(data.WithID(id))
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return data.ZeroTask, fmt.Errorf("%w: %v", data.ErrIDNotFound, err)
		}
		return data.ZeroTask, err
	}
	if len(tasks) != 1 {
		return data.ZeroTask, data.ErrAmbiguousID
	}
	return tasks[0], nil
}

func (c *TaskController) GetAllTasks() ([]data.Task, error) {
	tasks, err := c.store.Get()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *TaskController) Count() (int, error) {
	return c.store.Count()
}

func (c *TaskController) CreateTask(task data.Task) (data.Task, error) {
	created, err := c.store.Create(task)
	if err != nil {
		return data.ZeroTask, err
	}
	return created, nil
}

func (c *TaskController) UpdateTask(id uint64, task data.Task) (data.Task, error) {
	updated, err := c.store.Update(id, task)
	if err != nil {
		return data.ZeroTask, err
	}
	return updated, nil
}

func (c *TaskController) DeleteTask(id uint64) error {
	if err := c.store.Delete(id); err != nil {
		return err
	}
	return nil
}
