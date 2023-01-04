package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"sync"
	"time"
	"todo-list/data"
	"todo-list/eventlog"
)

const (
	TaskStartedEvent = "taskStarted"
	TaskStoppedEvent = "taskStopped"
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

func (c *TaskController) onStartup(ctx context.Context) error {
	c.ctx = ctx
	entry, err := c.store.GetRunningTimeEntry()
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}
	c.activeTimeEntry = entry
	runtime.EventsEmit(ctx, TaskStartedEvent, entry.TaskID)
	return nil
}

func (c *TaskController) GetTimeEntriesToday() ([]data.TimeEntry, error) {
	return c.store.GetTimeEntries(data.EntriesToday())
}

func (c *TaskController) GetTimeEntriesForWeek() ([]data.TimeEntry, error) {
	return c.store.GetTimeEntries(data.SinceWeekday(time.Sunday))
}

type jsObject = map[string]any

func stoppedEvent(id uint64) jsObject {
	return jsObject{
		"stopped": id,
	}
}

func startedEvent(id uint64) jsObject {
	return jsObject{
		"started": id,
	}
}

type TrackedTaskDetails struct {
	Task  data.Task      `json:"task"`
	Entry data.TimeEntry `json:"entry"`
}

func (c *TaskController) GetTrackedTaskDetails() (*TrackedTaskDetails, error) {
	if c.activeTimeEntry == nil {
		return nil, nil
	}
	tasks, err := c.store.Get(data.WithID(c.activeTimeEntry.TaskID))
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("unable to locate task by ID %d", c.activeTimeEntry.TaskID)
	}
	return &TrackedTaskDetails{
		Task:  tasks[0],
		Entry: *c.activeTimeEntry,
	}, nil
}

func (c *TaskController) StartTask(taskID uint64) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.activeTimeEntry != nil {
		oldID := c.activeTimeEntry.TaskID
		newEntry, err := c.store.StartAfterStop(taskID, c.activeTimeEntry.ID)
		if err != nil {
			return err
		}
		c.activeTimeEntry = &newEntry
		runtime.EventsEmit(c.ctx, TaskStoppedEvent, stoppedEvent(oldID))
		runtime.EventsEmit(c.ctx, TaskStartedEvent, startedEvent(newEntry.TaskID))
		return nil
	}

	entry, err := c.store.StartTimeEntry(taskID)
	if err != nil {
		return err
	}
	c.activeTimeEntry = &entry
	runtime.EventsEmit(c.ctx, TaskStartedEvent, startedEvent(entry.TaskID))
	return nil
}

func (c *TaskController) StopTask() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.activeTimeEntry == nil {
		return nil
	}
	taskID := c.activeTimeEntry.TaskID
	err := c.store.StopTimeEntry(c.activeTimeEntry.ID)
	if err != nil {
		return err
	}

	c.activeTimeEntry = nil
	runtime.EventsEmit(c.ctx, TaskStoppedEvent, stoppedEvent(taskID))
	return err
}

func (c *TaskController) GetTaskByID(id uint64) (data.Task, error) {
	c.log.DebugEvent("Received GetTasksByID call", "id", id)
	tasks, err := c.store.Get(data.WithID(id))
	if err != nil {
		if errors.Is(err, bolthold.ErrNotFound) {
			return data.Task{}, fmt.Errorf("%w: %v", data.ErrIDNotFound, err)
		}
		return data.Task{}, err
	}
	if len(tasks) != 1 {
		return data.Task{}, data.ErrAmbiguousID
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
		return data.Task{}, err
	}
	return created, nil
}

func (c *TaskController) UpdateTask(id uint64, task data.Task) (data.Task, error) {
	c.log.DebugEvent("Received UpdateTask call", "id", id, "task", task)
	updated, err := c.store.Update(id, task)
	if err != nil {
		return data.Task{}, err
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
