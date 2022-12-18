package data

import "github.com/timshannon/bolthold"

type TaskFilter = func(query *bolthold.Query)

func WithID(id uint64) TaskFilter {
	return func(query *bolthold.Query) {
		query.And(bolthold.Key).Eq(id)
	}
}

type MergeStrategy = string

const (
	MergeKeepInternal MergeStrategy = "KeepInternal" // Refuse conflicting updates.
	MergeOverwrite    MergeStrategy = "Overwrite"    // Overwrite the local store with data from the snapshot.
	MergeError        MergeStrategy = "Error"        // Don't attempt to merge, just return an error.
)

// TaskStorage facilitates persistence operations for Task data.
type TaskStorage interface {
	// Get will retrieve Tasks that match the given filter(s), or all Tasks.
	Get(...TaskFilter) ([]Task, error)
	// Count will count the number of Tasks in the data store.
	Count() (int, error)
	// Create will create a new Task using the given template Task.
	Create(template Task) (Task, error)
	// Update will update the state of a Task referenced by id with the given template Task.
	Update(id uint64, template Task) (Task, error)
	// Delete will remove a Task record from the store.
	Delete(id uint64) error
	// Export will write all Task data in the store to a timestamped CSV file in the given directory.
	// If Export is successful, it will return the file system location of the exported data.
	Export(dir string) (string, error)
	// Import will read all Task data from the given file into the store.
	// If a Task read from the file conflicts with the state of the store, then the given merge strategy will be used.
	Import(file string, merge MergeStrategy) error
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
