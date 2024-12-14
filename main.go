package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Details  string `json:"details"`
	UserName string `json:"user_name,omitempty"`
}

var users = []User{
	{ID: 1, Name: "John Doe"},
	{ID: 2, Name: "Jane Smith"},
}

var orders = []Order{}

func main() {
	go func() {
		router := gin.Default()
		router.GET("/users/:id", getUserByID)
		router.POST("/users", createUser)
		log.Println("ttp://localhost:8081")
		router.Run(":8081")
	}()

	router := gin.Default()
	router.POST("/orders", createOrder)
	router.GET("/orders", getOrders)
	log.Println("http://localhost:8082")
	router.Run(":8082")
}

func getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, user)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func createUser(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users = append(users, newUser)
	c.JSON(http.StatusCreated, newUser)
}

func createOrder(c *gin.Context) {
	var newOrder Order
	if err := c.ShouldBindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Fetching user with ID: %d", newOrder.UserID)
	user, err := fetchUserByID(newOrder.UserID)
	if err != nil {
		log.Printf("Error fetching user: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	newOrder.UserName = user.Name
	orders = append(orders, newOrder)
	c.JSON(http.StatusCreated, newOrder)
}

func getOrders(c *gin.Context) {
	c.JSON(http.StatusOK, orders)
}

func fetchUserByID(userID int) (*User, error) {
	client := resty.New()
	url := fmt.Sprintf("http://localhost:8081/users/%d", userID)
	resp, err := client.R().
		SetResult(&User{}).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("user not found")
	}
	return resp.Result().(*User), nil
}
