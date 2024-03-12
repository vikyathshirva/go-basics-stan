package user

import (
	"fmt"
	"sync"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserRepository interface {
	Create(user *User) error
	Read(userID int) (*User, error)
	Update(user *User) error
	Delete(userID int) error
}

type InMemoryUserRepository struct {
	users   map[int]*User
	counter int
	mu      sync.Mutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:   make(map[int]*User),
		counter: 1,
	}
}

func (r *InMemoryUserRepository) Create(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	user.ID = r.counter
	r.users[user.ID] = user
	r.counter++
	return nil
}

func (r *InMemoryUserRepository) Read(userID int) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) Update(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.users[user.ID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(userID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.users[userID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	delete(r.users, userID)
	return nil
}
