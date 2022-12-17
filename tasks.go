package main

import (
	"github.com/timshannon/bolthold"
	"sort"
)

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

type TaskController struct {
	store TaskStorage
}

var _ sort.Interface = (*taskSorter)(nil)

type taskSorter struct {
	tasks []Task
}

func newTaskSorter(tasks []Task) *taskSorter {
	return &taskSorter{tasks: tasks}
}

func (t *taskSorter) Len() int {
	return len(t.tasks)
}

func (t *taskSorter) Less(i, j int) bool {
	a, b := t.tasks[i], t.tasks[j]

	switch {
	case a.Done && !b.Done:
		return false
	case !a.Done && b.Done:
		return true
	case a.ID < b.ID:
		return true
	default:
		return false
	}
}

func (t *taskSorter) Swap(i, j int) {
	t.tasks[i], t.tasks[j] = t.tasks[j], t.tasks[i]
}

func NewTaskController() (*TaskController, error) {
	storage, err := newBoltStorage()
	if err != nil {
		return nil, err
	}
	return &TaskController{
		store: storage,
	}, nil
}

func (c *TaskController) GetTaskByID(id uint64) (Task, error) {
	tasks, err := c.store.Get(WithID(id))
	if err != nil {
		return zeroTask, err
	}
	if len(tasks) != 1 {
		return zeroTask, ErrAmbiguousID
	}
	return tasks[0], nil
}

func (c *TaskController) GetAllTasks() ([]Task, error) {
	tasks, err := c.store.Get()
	if err != nil {
		return nil, err
	}
	sort.Sort(newTaskSorter(tasks))
	return tasks, nil
}

func (c *TaskController) Count() (int, error) {
	return c.store.Count()
}

func (c *TaskController) CreateTask(task Task) (Task, error) {
	created, err := c.store.Create(task)
	if err != nil {
		return zeroTask, err
	}
	return created, nil
}

func (c *TaskController) UpdateTask(id uint64, task Task) (Task, error) {
	updated, err := c.store.Update(id, task)
	if err != nil {
		return zeroTask, err
	}
	return updated, nil
}

func (c *TaskController) DeleteTask(id uint64) error {
	if err := c.store.Delete(id); err != nil {
		return err
	}
	return nil
}
