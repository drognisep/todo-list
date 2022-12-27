package data

import (
	"time"
)

type TimeEntry struct {
	ID     uint64     `json:"id" boltholdKey:"ID"`
	TaskID uint64     `json:"taskID" boltholdIndex:"TaskID"`
	Start  time.Time  `json:"start" boltholdIndex:"Start"`
	End    *time.Time `json:"end"`
	Synced bool       `json:"synced" boltholdIndex:"Synced"`
}

func (e *TimeEntry) Duration() time.Duration {
	if e.End == nil {
		return time.Duration(0)
	}
	return e.End.Sub(e.Start)
}
