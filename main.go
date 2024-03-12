package main

import (
	"net/http"
	"strconv"
	"github.com/vikyathshirva/go-basics-stan/user"
	"github.com/gin-gonic/gin"
)

func setupRoutes(userRepo user.UserRepository) *gin.Engine {
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

	var user user.User
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
	inMemoryRepo := user.NewInMemoryUserRepository()

	go func() {
		memRouter := setupRoutes(inMemoryRepo)
		memRouter.Run(":8080")
	}()

	select {}
}

