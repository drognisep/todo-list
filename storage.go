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

func refSequence(seqStart ...uint64) func() uint64 {
	var mux sync.Mutex
	next := uint64(1)
	if len(seqStart) == 1 {
		next = seqStart[0]
	}

	return func() uint64 {
		mux.Lock()
		defer mux.Unlock()
		val := next
		next++
		return val
	}
}

var _ TaskStorage = (*taskStorage)(nil)

type taskStorage struct {
	mux   sync.RWMutex
	id    []uint64
	name  []string
	desc  []string
	done  []bool
	idIdx index[uint64]
	seq   func() uint64
}

func newTaskStorage() *taskStorage {
	return &taskStorage{
		idIdx: newIndex[uint64](),
		seq:   refSequence(),
	}
}

func (t *taskStorage) all() refSet {
	set := refSet{}
	for i := range t.idIdx.ref {
		set.add(i)
	}
	return set
}

func (t *taskStorage) Get(filters ...TaskFilter) ([]Task, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	var tasks []Task
	refs := t.all()
collect:
	for ref := range refs {
		for _, filter := range filters {
			if !filter(t, ref) {
				continue collect
			}
		}
		tasks = append(tasks, Task{
			ID:          t.id[ref],
			Name:        t.name[ref],
			Description: t.desc[ref],
			Done:        t.done[ref],
		})
	}
	return tasks, nil
}

func (t *taskStorage) Count() (int, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return lenIndex(t.idIdx), nil
}

func (t *taskStorage) getRef(ref int) (Task, bool) {
	id, ok := t.idIdx.ref[ref]
	if !ok {
		return zeroTask, false
	}

	return Task{
		ID:          id,
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
	newRef := len(t.id)
	newID := t.seq()
	task.ID = newID
	t.id = append(t.id, newID)
	t.name = append(t.name, task.Name)
	t.desc = append(t.desc, task.Description)
	t.done = append(t.done, task.Done)
	addIndex(t.idIdx, newRef, newID)

	return task, nil
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
	id := newState.ID
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
