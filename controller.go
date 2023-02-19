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
	TaskStartedEvent    = "taskStarted"
	TaskStoppedEvent    = "taskStopped"
	NotesUpdatedEvent   = "notesUpdated"
	EntriesUpdatedEvent = "entriesChanged"
)

type ModelController struct {
	ctx   context.Context
	store data.TaskStorage
	log   *eventlog.EventLog

	mux             sync.Mutex
	activeTimeEntry *data.TimeEntry
}

func NewTaskController(logger *eventlog.EventLog) (*ModelController, error) {
	store, err := data.NewStorage()
	if err != nil {
		return nil, err
	}
	return &ModelController{
		store: store,
		log:   logger,
	}, nil
}

func (c *ModelController) onStartup(ctx context.Context) error {
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

func (c *ModelController) GetTimeEntriesToday() ([]data.TimeEntry, error) {
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

func (c *ModelController) GetTimeEntriesForWeek() ([]data.TimeEntry, error) {
	return c.store.GetTimeEntries(data.SinceWeekday(time.Sunday))
}

func (c *ModelController) GetSummaryForEntries(entries []data.TimeEntry) (*TimeEntrySummary, error) {
	return c.calculateSummary(entries)
}

func (c *ModelController) calculateSummary(entries []data.TimeEntry) (*TimeEntrySummary, error) {
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
			tasks, err := c.store.GetHistoricTasks(data.WithID(e.TaskID))
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

func (c *ModelController) GetTrackedTaskDetails() (*TrackedTaskDetails, error) {
	if c.activeTimeEntry == nil {
		return nil, nil
	}
	tasks, err := c.store.GetTasks(data.WithID(c.activeTimeEntry.TaskID))
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

func (c *ModelController) StartTask(taskID uint64) error {
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

func (c *ModelController) StopTask() error {
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

func (c *ModelController) GetTaskByID(id uint64) (data.Task, error) {
	c.log.DebugEvent("Received GetTasksByID call", "id", id)
	tasks, err := c.store.GetHistoricTasks(data.WithID(id))
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

func (c *ModelController) GetAllTasks() ([]data.Task, error) {
	c.log.DebugEvent("Received GetAllTasks call")
	tasks, err := c.store.GetTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *ModelController) Count() (int, error) {
	c.log.DebugEvent("Received CountTasks call")
	return c.store.CountTasks()
}

func (c *ModelController) CreateTask(task data.Task) (data.Task, error) {
	c.log.DebugEvent("Received CreateTask call", "toCreate", task)
	created, err := c.store.CreateTask(task)
	if err != nil {
		return data.Task{}, err
	}
	return created, nil
}

func (c *ModelController) UpdateTask(id uint64, task data.Task) (data.Task, error) {
	c.log.DebugEvent("Received UpdateTask call", "id", id, "task", task)
	updated, err := c.store.UpdateTask(id, task)
	if err != nil {
		return data.Task{}, err
	}
	return updated, nil
}

func (c *ModelController) DeleteTask(id uint64) error {
	c.log.DebugEvent("Received DeleteTask call", "id", id)
	if err := c.store.DeleteTask(id); err != nil {
		return err
	}
	return nil
}

func (c *ModelController) GetTaskNotes(taskID uint64) ([]data.Note, error) {
	c.log.DebugEvent("Received GetTaskNotes call", "id", taskID)
	return c.store.GetTaskNotes(taskID)
}

func (c *ModelController) GetTaskNoteCount(taskID uint64) (int, error) {
	c.log.DebugEvent("Received GetTaskNoteCount call", "id", taskID)
	return c.store.GetTaskNoteCount(taskID)
}

func (c *ModelController) AddNote(taskID uint64, note string) (data.Note, error) {
	c.log.DebugEvent("Received AddNote call", "id", taskID, "noteText", note)
	addNote, err := c.store.AddNote(taskID, data.Note{Text: note})
	if err == nil {
		runtime.EventsEmit(c.ctx, NotesUpdatedEvent)
	}
	return addNote, err
}

func (c *ModelController) UpdateNote(noteID uint64, note string) (data.Note, error) {
	c.log.DebugEvent("Received UpdateNote call", "id", noteID, "newText", note)
	updateNote, err := c.store.UpdateNote(noteID, data.Note{Text: note})
	if err == nil {
		runtime.EventsEmit(c.ctx, NotesUpdatedEvent)
	}
	return updateNote, err
}

func (c *ModelController) DeleteNote(noteID uint64) error {
	c.log.DebugEvent("Received DeleteNote call", "id", noteID)
	err := c.store.DeleteNote(noteID)
	if err == nil {
		runtime.EventsEmit(c.ctx, NotesUpdatedEvent)
	}
	return err
}

func (c *ModelController) Export() (string, error) {
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

func (c *ModelController) Import(strategy string) error {
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
