package data

import (
	"errors"
	"github.com/timshannon/bolthold"
	"time"
)

var (
	ErrIDNotFound              = errors.New("specified ID not found")
	ErrAmbiguousID             = errors.New("ambiguous ID detected")
	ErrUnmappedReqdImportField = errors.New("unable to map required import field")
)

type TaskFilter = func(query *bolthold.Query)

func WithID(id uint64) TaskFilter {
	return func(query *bolthold.Query) {
		query.And(bolthold.Key).Eq(id).Limit(1)
	}
}

type TimeEntryFilter = func(query *bolthold.Query)

// After will only return TimeEntry records that start on or after the given time.Time.
func After(after time.Time) TimeEntryFilter {
	return func(query *bolthold.Query) {
		query.And("Start").Ge(after)
	}
}

// Before will only return TimeEntry records that start on or before the given time.Time.
func Before(before time.Time) TimeEntryFilter {
	return func(query *bolthold.Query) {
		query.And("End").Le(before)
	}
}

// NotSynced will only return TimeEntry records that have not been as synced to an external store.
func NotSynced() TimeEntryFilter {
	return func(query *bolthold.Query) {
		query.And("Synced").Eq(false)
	}
}

// ForTask will only return TimeEntry records for a Task identified by its ID.
func ForTask(taskID uint64) TimeEntryFilter {
	return func(query *bolthold.Query) {
		query.And("TaskID").Eq(taskID)
	}
}

// EntriesToday will return an After filter for TimeEntry records that occurred today.
func EntriesToday() TimeEntryFilter {
	now := time.Now()
	after := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return After(after)
}

// SinceWeekday will return an After filter for the previous occurrence of the given weekday, at the beginning of the day.
// If the goal time.Weekday matches that of the current day, then the filter will capture all TimeEntry records from 7 days ago.
func SinceWeekday(goal time.Weekday) TimeEntryFilter {
	return After(lastWeekday(goal, time.Now()))
}

func lastWeekday(goal time.Weekday, given time.Time) time.Time {
	year, month, day := given.Date()
	weekday := given.Weekday()
	offset := int(goal - weekday)
	if offset >= 0 {
		offset -= 7
	}
	return time.Date(year, month, day+offset, 0, 0, 0, 0, given.Location())
}

// TaskStorage facilitates persistence operations for Task data.
type TaskStorage interface {
	// GetTasks will retrieve Tasks that match the given filter(s), or all Tasks.
	GetTasks(...TaskFilter) ([]Task, error)
	// GetHistoricTasks will behave like Get, except that it will also get soft-deleted Tasks.
	GetHistoricTasks(...TaskFilter) ([]Task, error)
	// CountTasks will count the number of Tasks in the data store.
	CountTasks() (int, error)
	// CreateTask will create a new Task using the given template Task.
	CreateTask(template Task) (Task, error)
	// UpdateTask will update the state of a Task referenced by id with the given template Task.
	UpdateTask(id uint64, template Task) (Task, error)
	// DeleteTask will inactivate a Task record in the store by marking it as soft-deleted.
	DeleteTask(id uint64) error
	// Export will write all Task data in the store to a timestamped CSV file in the given directory.
	// If Export is successful, it will return the file system location of the exported data.
	Export(dir string) (string, error)
	// Import will read all Task data from the given file into the store.
	// If a Task read from the file conflicts with the state of the store, then the given merge strategy will be used.
	Import(file string, merge MergeStrategy) error
	// StartTimeEntry will start progress on the Task identified by taskID.
	// Returns a new TimeEntry when successful.
	StartTimeEntry(taskID uint64) (TimeEntry, error)
	// StopTimeEntry stops tracking for the given TimeEntry.
	// Returns an error if the operation was not able to be completed.
	StopTimeEntry(entryID uint64) error
	// StartAfterStop will start a new TimeEntry with startTaskID and stop the current entry in a single transaction.
	// The Task's start time will match the TimeEntry's end time.
	StartAfterStop(startTaskID uint64, stopEntryID uint64) (TimeEntry, error)
	// GetTimeEntries will return TimeEntry records that match the given criteria, or all if none are given.
	GetTimeEntries(filters ...TimeEntryFilter) ([]TimeEntry, error)
	// GetRunningTimeEntry returns the latest running time entry, if one exist.
	GetRunningTimeEntry() (*TimeEntry, error)
	// AddNote adds a Note to a Task identified by taskID.
	AddNote(taskID uint64, note Note) (Note, error)
	// GetTaskNotes gets a list of Note related to the Task identified by taskID.
	GetTaskNotes(taskID uint64) ([]Note, error)
	// GetTaskNoteCount gets the number of Note related to the TAsk identified by taskID.
	GetTaskNoteCount(taskID uint64) (int, error)
	// UpdateNote updates the text of the Note identified by noteID.
	UpdateNote(noteID uint64, note Note) (Note, error)
	// DeleteNote deletes a Note identified by noteID.
	DeleteNote(noteID uint64) error
}

// Storage combines the existing persistence interfaces into one for convenience.
type Storage interface {
	TaskStorage
}

// NewStorage creates or loads the store from the default location.
func NewStorage() (Storage, error) {
	return newBoltStorage()
}

// NewStorageAt creates or loads the store at the given location.
func NewStorageAt(file string) (Storage, error) {
	return newBoltStorageAt(file)
}
