package main

import (
	"fmt"
	"net/http"
	"sync"
	"strconv"
	"github.com/gin-gonic/gin"
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

type DBUserRepository struct {
}

func setupRoutes(userRepo UserRepository) *gin.Engine {
	router := gin.Default()

	router.GET("/user/:id", func(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := userRepo.Read(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
})

router.PUT("/user/:id", func(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user.ID = id
	if err := userRepo.Update(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
})

router.DELETE("/user/:id", func(c *gin.Context) {
	userID := c.Param("id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := userRepo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
})


	return router
}

func main() {
	inMemoryRepo := NewInMemoryUserRepository()

	go func() {
		memRouter := setupRoutes(inMemoryRepo)
		memRouter.Run(":8080")
	}()

	select {}
}
