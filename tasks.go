package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"sync"
	"todo-list/data"
	"todo-list/eventlog"
)

type TaskController struct {
	ctx   context.Context
	store data.TaskStorage
	log   *eventlog.EventLog

	mux             sync.Mutex
	activeTimeEntry *data.TimeEntry
}

func NewTaskController(logger *eventlog.EventLog) (*TaskController, error) {
	store, err := data.NewStorage()
	if err != nil {
		return nil, err
	}
	return &TaskController{
		store: store,
		log:   logger,
	}, nil
}

func (c *TaskController) StartTask(taskID uint64) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, err := c.store.StartTimeEntry(taskID)
	if err != nil {
		return err
	}
	c.activeTimeEntry = &entry
	return nil
}

func (c *TaskController) StopTask() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.activeTimeEntry == nil {
		return nil
	}
	err := c.store.StopTimeEntry(*c.activeTimeEntry)
	if err != nil {
		return err
	}

	c.activeTimeEntry = nil
	return err
}

func (c *TaskController) GetTaskByID(id uint64) (data.Task, error) {
	c.log.DebugEvent("Received GetTasksByID call", "id", id)
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
	c.log.DebugEvent("Received GetAllTasks call")
	tasks, err := c.store.Get()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *TaskController) Count() (int, error) {
	c.log.DebugEvent("Received Count call")
	return c.store.Count()
}

func (c *TaskController) CreateTask(task data.Task) (data.Task, error) {
	c.log.DebugEvent("Received CreateTask call", "toCreate", task)
	created, err := c.store.Create(task)
	if err != nil {
		return data.ZeroTask, err
	}
	return created, nil
}

func (c *TaskController) UpdateTask(id uint64, task data.Task) (data.Task, error) {
	c.log.DebugEvent("Received UpdateTask call", "id", id, "task", task)
	updated, err := c.store.Update(id, task)
	if err != nil {
		return data.ZeroTask, err
	}
	return updated, nil
}

func (c *TaskController) DeleteTask(id uint64) error {
	c.log.DebugEvent("Received DeleteTask call", "id", id)
	if err := c.store.Delete(id); err != nil {
		return err
	}
	return nil
}

func (c *TaskController) Export() (string, error) {
	c.log.DebugEvent("Received Export call")
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir, err := runtime.OpenDirectoryDialog(c.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:     home,
		CanCreateDirectories: true,
	})
	if err != nil {
		return "", err
	}

	if len(dir) == 0 {
		return "", nil
	}

	return c.store.Export(dir)
}

func (c *TaskController) Import(strategy string) error {
	c.log.DebugEvent("Received Import call")
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file, err := runtime.OpenFileDialog(c.ctx, runtime.OpenDialogOptions{
		DefaultDirectory: home,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Snapshots",
				Pattern:     "*.snapshot",
			},
		},
	})
	if err != nil {
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return c.store.Import(file, strategy)
}
