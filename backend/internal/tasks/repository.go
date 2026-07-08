package tasks

import (
	"errors"
	"math/rand"
	"sync"
)

var ErrTaskNotFound = errors.New("task not found")

type Repository interface {
	Seed([]Task)
	RandomByCategory(Difficulty) (Task, error)
	FindByID(string) (Task, error)
}

type MemoryRepository struct {
	mu    sync.RWMutex
	tasks map[string]Task
	byCat map[Difficulty][]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{tasks: map[string]Task{}, byCat: map[Difficulty][]string{}}
}

func (r *MemoryRepository) Seed(seed []Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range seed {
		r.tasks[task.ID] = task
		r.byCat[task.Category] = append(r.byCat[task.Category], task.ID)
	}
}

func (r *MemoryRepository) RandomByCategory(category Difficulty) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := r.byCat[category]
	if len(ids) == 0 {
		return Task{}, ErrTaskNotFound
	}
	return r.tasks[ids[rand.Intn(len(ids))]], nil
}

func (r *MemoryRepository) FindByID(id string) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}
