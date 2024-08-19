package main

import (
	"log"
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func stringPtr(s string) *string {
	return &s
}

// Sample users with pointer values for optional fields

var users = []User{
	{
		ID:        "1",
		FirstName: stringPtr("John"),
		Surname:   stringPtr("Doe"),
		Email:     "john@example.com",
	},
	{
		ID:        "2",
		FirstName: stringPtr("Arthur"),
		Surname:   stringPtr("Morgan"),
		Email:     "arthur@example.com",
	},
	{
		ID:        "3",
		FirstName: stringPtr("Fake"),
		Surname:   stringPtr("Name"),
		Email:     "fake@example.com",
	},
}

var maxId string = "3"

func incrementMaxId() string {
	intVal, _ := (strconv.Atoi(maxId))

	intVal++

	maxId = strconv.Itoa(intVal)
	return maxId
}

func postUser(c *gin.Context) {
	var userPost UserPost
	if err := c.BindJSON(&userPost); err != nil {
		return
	}

	//salt and hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPost.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	//increment id and then add
	newUser := User{incrementMaxId(), userPost.FirstName, userPost.Surname, userPost.Email, hashedPassword}
	// Add the new album to the slice.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserById(c *gin.Context) {

	id := c.Param("id")

	for _, a := range users {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func main() {
	router := gin.Default()
	router.GET("/:id", getUserById)
	router.POST("/", postUser)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Users Microservice"})
	})

	router.Run(":8080")
}
