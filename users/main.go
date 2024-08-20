package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	gin "github.com/gin-gonic/gin"
	bcrypt "golang.org/x/crypto/bcrypt"
)

func stringPtr(s string) *string {
	return &s
}

var jwtKey = []byte("my_secret_key")
var logger *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

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

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user *User = nil
	for _, u := range users {
		if u.Email == creds.Username {
			user = &u
		}
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Validate credentials here (usually check against a database)
	err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(creds.Password))
	if err != nil {
		// Password does not match
		logger.Println("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	logger.Println("Password is correct")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(30 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func protectedEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "This is a protected endpoint"})
}

func main() {
	router := gin.Default()
	router.GET("users/:id", getUserById)
	router.POST("users/", postUser)

	router.GET("users/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Users Microservice"})
	})

	router.POST("users/login", login)

	router.GET("users/protected", AuthMiddleware(), protectedEndpoint)

	router.Run(":8080")
}
