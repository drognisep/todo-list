package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/timshannon/bolthold"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"sort"
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

var _ sort.Interface = (*TimeEntrySummary)(nil)

type TimeEntrySummary struct {
	Lines []TaskSummary `json:"lines"`
	Total string        `json:"total"`
}

func (t *TimeEntrySummary) Len() int {
	return len(t.Lines)
}

func (t *TimeEntrySummary) Less(i, j int) bool {
	return t.Lines[i].totalDur < t.Lines[j].totalDur
}

func (t *TimeEntrySummary) Swap(i, j int) {
	t.Lines[i], t.Lines[j] = t.Lines[j], t.Lines[i]
}

type TaskSummary struct {
	Name     string `json:"name"`
	Total    string `json:"duration"`
	totalDur time.Duration
}

func (c *TaskController) GetTimeEntriesForWeek() ([]data.TimeEntry, error) {
	return c.store.GetTimeEntries(data.SinceWeekday(time.Sunday))
}

func (c *TaskController) GetSummaryForEntries(entries []data.TimeEntry) (*TimeEntrySummary, error) {
	return c.calculateSummary(entries)
}

func (c *TaskController) calculateSummary(entries []data.TimeEntry) (*TimeEntrySummary, error) {
	if len(entries) == 0 {
		return &TimeEntrySummary{}, nil
	}

	var (
		taskSummary = map[uint64]TaskSummary{}
		entryTotals = map[uint64]time.Duration{}
		total       = time.Duration(0)
		lines       []TaskSummary
	)

	for _, e := range entries {
		if e.End == nil {
			continue
		}
		if _, ok := taskSummary[e.TaskID]; !ok {
			tasks, err := c.store.GetHistoric(data.WithID(e.TaskID))
			if err != nil {
				return nil, err
			}
			if len(tasks) == 0 {
				return nil, fmt.Errorf("%w: unable to locate task ID '%d' in entry ID '%d'", data.ErrIDNotFound, e.TaskID, e.ID)
			}
			taskSummary[e.TaskID] = TaskSummary{Name: tasks[0].Name}
		}

		dur := e.End.Sub(e.Start).Round(time.Second)
		total += dur
		var runningTotal time.Duration

		if cachedTotal, ok := entryTotals[e.TaskID]; ok {
			runningTotal = cachedTotal
		}
		entryTotals[e.TaskID] = runningTotal + dur
	}

	for tid, taskTotal := range entryTotals {
		line := taskSummary[tid]
		line.totalDur = taskTotal
		line.Total = taskTotal.String()
		lines = append(lines, line)
	}

	summary := &TimeEntrySummary{Lines: lines, Total: total.String()}
	sort.Sort(sort.Reverse(summary))
	return summary, nil
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
	tasks, err := c.store.GetHistoric(data.WithID(id))
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
