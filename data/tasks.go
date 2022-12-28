package data

type Task struct {
	ID          uint64 `json:"id" boltholdKey:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Done        bool   `json:"done" boltholdIndex:"Done"`
	Priority    int    `json:"priority" boltholdIndex:"Priority"`
	Favorite    bool   `json:"favorite"`
	SoftDeleted bool   `json:"inactivated" boltholdIndex:"SoftDeleted"`
}

type exportModel struct {
	Tasks       []Task      `json:"tasks"`
	TimeEntries []TimeEntry `json:"timeEntries,omitempty"`
}
