package main

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrIDNotFound = errors.New("specified ID not found")
	zeroTask      = Task{}
)

type index[T comparable] struct {
	val map[T]int
	ref map[int]T
}

func newIndex[T comparable]() index[T] {
	return index[T]{
		val: map[T]int{},
		ref: map[int]T{},
	}
}

func addIndex[T comparable](idx index[T], i int, t T) {
	idx.val[t] = i
	idx.ref[i] = t
}

func delIndex[T comparable](idx index[T], i int) {
	t, ok := idx.ref[i]
	if !ok {
		return
	}
	delete(idx.ref, i)
	delete(idx.val, t)
}

func lenIndex[T comparable](idx index[T]) int {
	return len(idx.val)
}

type refSet map[int]bool

func (s refSet) add(i int) {
	s[i] = true
}

func (s refSet) has(i int) bool {
	return s[i]
}

func (s refSet) intersect(other refSet) refSet {
	if s == nil {
		return refSet{}
	}

	xref := refSet{}
	for idx := range other {
		if s.has(idx) {
			xref.add(idx)
		}
	}
	return xref
}

var taskSeq = func() func() uint64 {
	var mux sync.Mutex
	start := uint64(1)
	return func() uint64 {
		mux.Lock()
		defer mux.Unlock()
		rval := start
		start++
		return rval
	}
}()

var _ TaskStorage = (*taskStorage)(nil)

type taskStorage struct {
	mux   sync.RWMutex
	id    []uint64
	name  []string
	desc  []string
	done  []bool
	idIdx index[uint64]
}

func newTaskStorage() *taskStorage {
	return &taskStorage{
		idIdx: newIndex[uint64](),
	}
}

func (t *taskStorage) all() refSet {
	set := refSet{}
	for i := range t.idIdx.ref {
		set.add(i)
	}
	return set
}

func (t *taskStorage) Get() ([]Task, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	tblLen := lenIndex(t.idIdx)
	tasks := make([]Task, tblLen)
	refs := t.all()
	for i := range refs {
		tasks[i] = Task{
			TaskID:      t.id[i],
			Name:        t.name[i],
			Description: t.desc[i],
			Done:        t.done[i],
		}
	}
	return tasks, nil
}

func (t *taskStorage) getRef(ref int) (Task, bool) {
	id, ok := t.idIdx.ref[ref]
	if !ok {
		return zeroTask, false
	}

	return Task{
		TaskID:      id,
		Name:        t.name[ref],
		Description: t.desc[ref],
		Done:        t.done[ref],
	}, true
}

func (t *taskStorage) getID(id uint64) (Task, bool) {
	ref, ok := t.idIdx.val[id]
	if !ok {
		return zeroTask, false
	}
	return t.getRef(ref)
}

func (t *taskStorage) Create(task Task) (Task, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	newTask := task.Copy()
	newRef := len(t.id)
	newID := taskSeq()
	newTask.TaskID = newID
	t.id = append(t.id, newID)
	t.name = append(t.name, task.Name)
	t.desc = append(t.desc, task.Description)
	t.done = append(t.done, task.Done)
	addIndex(t.idIdx, newRef, newID)

	return newTask, nil
}

func (t *taskStorage) Update(id uint64, req Task) (Task, error) {
	t.mux.Lock()
	defer t.mux.Unlock()
	_new, ok := t.getID(id)
	if !ok {
		return zeroTask, fmt.Errorf("%w: %d", ErrIDNotFound, id)
	}

	_new.Name = req.Name
	_new.Description = req.Description
	_new.Done = req.Done

	return _new, t.update(_new)
}

func (t *taskStorage) update(newState Task) error {
	id := newState.TaskID
	ref, _ := t.idIdx.val[id]

	t.id[ref] = id
	t.name[ref] = newState.Name
	t.desc[ref] = newState.Description
	t.done[ref] = newState.Done
	return nil
}

func (t *taskStorage) Delete(id uint64) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	ref, ok := t.idIdx.val[id]
	if !ok {
		return fmt.Errorf("%w: %d", ErrIDNotFound, id)
	}
	delIndex(t.idIdx, ref)
	return nil
}
