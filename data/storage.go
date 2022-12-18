package data

import "github.com/timshannon/bolthold"

type TaskFilter = func(query *bolthold.Query)

func WithID(id uint64) TaskFilter {
	return func(query *bolthold.Query) {
		query.And(bolthold.Key).Eq(id)
	}
}

type TaskStorage interface {
	Get(...TaskFilter) ([]Task, error)
	Count() (int, error)
	Create(Task) (Task, error)
	Update(uint64, Task) (Task, error)
	Delete(uint64) error
}

type Storage interface {
	TaskStorage
}

func NewStorage() (Storage, error) {
	return newBoltStorage()
}

func NewStorageAt(file string) (Storage, error) {
	return newBoltStorageAt(file)
}
