package main

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

func NewTaskController() (*TaskController, error) {
	return &TaskController{
		store: newTaskStorage(),
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
	return tasks, nil
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
