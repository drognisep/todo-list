package main

import (
	"errors"
)

type TaskFilter = func(t *taskStorage, ref int) bool

func WithID(id uint64) TaskFilter {
	return func(t *taskStorage, ref int) bool {
		return t.id[ref] == id
	}
}

type TaskStorage interface {
	Get(...TaskFilter) ([]Task, error)
	Count() (int, error)
	Create(Task) (Task, error)
	Update(uint64, Task) (Task, error)
	Delete(uint64) error
}

type TaskController struct {
	store TaskStorage
}

func NewTaskController(storage TaskStorage) (*TaskController, error) {
	if storage == nil {
		return nil, errors.New("nil TaskStorage dependency")
	}
	return &TaskController{
		store: storage,
	}, nil
}
