package users

import (
	"errors"
	"sync"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

type Repository interface {
	Create(User) error
	FindByEmail(string) (User, error)
	FindByID(string) (User, error)
	Update(User) error
	List() []User
}

type MemoryRepository struct {
	mu      sync.RWMutex
	byID    map[string]User
	byEmail map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{byID: map[string]User{}, byEmail: map[string]string{}}
}

func (r *MemoryRepository) Create(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byEmail[user.Email]; exists {
		return ErrUserExists
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user.ID
	return nil
}

func (r *MemoryRepository) FindByEmail(email string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byEmail[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return r.byID[id], nil
}

func (r *MemoryRepository) FindByID(id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *MemoryRepository) Update(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[user.ID]; !ok {
		return ErrUserNotFound
	}
	r.byID[user.ID] = user
	r.byEmail[user.Email] = user.ID
	return nil
}

func (r *MemoryRepository) List() []User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]User, 0, len(r.byID))
	for _, user := range r.byID {
		result = append(result, user)
	}
	return result
}
