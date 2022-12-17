package main

type Task struct {
	ID          uint64 `json:"id" boltholdKey:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done" boltholdIndex:"Done"`
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
	}
}

func (t *Task) Map(updates *Task) {
	if updates == nil {
		return
	}

}
