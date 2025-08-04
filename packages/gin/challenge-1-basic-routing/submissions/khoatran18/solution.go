package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	// TODO: Create Gin router

	// TODO: Setup routes
	// GET /users - Get all users
	// GET /users/:id - Get user by ID
	// POST /users - Create new user
	// PUT /users/:id - Update user
	// DELETE /users/:id - Delete user
	// GET /users/search - Search users by name

	// TODO: Start server on port 8080
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	// TODO: Return all users
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    users,
		Message: "",
		Error:   "",
		Code:    http.StatusOK,
	})
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// TODO: Get user by ID
	// Handle invalid ID format
	// Return 404 if user not found
	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, _ := findUserByID(userIDInt)
	if user == nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "",
		Error:   "",
		Code:    http.StatusOK,
	})
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	// TODO: Parse JSON request body
	// Validate required fields
	// Add user to storage
	// Return created user
	var newUser User
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	err = validateUser(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "",
		Error:   "",
		Code:    http.StatusCreated,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Parse JSON request body
	// Find and update user
	// Return updated user
	userID := c.Param("id")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Code:    http.StatusBadRequest,
			Error:   err.Error(),
		})
		return
	}

	_, index := findUserByID(userIDInt)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Code:    http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	var userData User
	err = c.ShouldBindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Invalid Request Data",
		})
		return
	}
	err = validateUser(userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Invalid Request Data",
		})
		return
	}
	users[index] = userData

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("%s Updated", userData.Name),
		Error:   "",
		Code:    http.StatusOK,
		Data:    users[index],
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	// TODO: Get user ID from path
	// Find and remove user
	// Return success message
	userID := c.Param("id")
	UserIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
			Message: "Invalid ID Request",
		})
		return
	}
	_, index := findUserByID(UserIDInt)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	users = append(users[:index], users[index+1:]...)
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Delete User Successfully",
		Error:   "",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	// TODO: Get name query parameter
	// Filter users by name (case-insensitive)
	// Return matching users
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid Name Query",
			Code:    http.StatusBadRequest,
		})
		return
	}

	usersData := []User{}
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), strings.ToLower(name)) {
			usersData = append(usersData, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    usersData,
		Message: "Find Users By Name Successfully",
		Error:   "",
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	// TODO: Implement user lookup
	// Return user pointer and index, or nil and -1 if not found
	for i := range users {
		if users[i].ID == id {
			return &users[i], i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	// TODO: Implement validation
	// Check required fields: Name, Email
	// Validate email format (basic check)
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Email == "" {
		return errors.New("Email is required")
	}
	return nil
}
