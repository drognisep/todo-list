package main

type Task struct {
	ID          uint64 `json:"id" boltholdKey:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done" boltholdIndex:"Done"`
	Priority    int    `json:"priority" boltholdIndex:"Priority"`
}

func (t *Task) Copy() Task {
	if t == nil {
		return Task{}
	}
	return Task{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Done:        t.Done,
		Priority:    t.Priority,
	}
}
