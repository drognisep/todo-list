package main

import (
	"errors"
)

type TaskStorage interface {
	Get() ([]Task, error)
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
