package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	gin "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	bcrypt "golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")
var logger *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

func postUser(c *gin.Context) {
	var userPost UserPost
	if err := c.BindJSON(&userPost); err != nil {
		return
	}

	//salt and hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPost.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//increment id and then add
	newUser := User{0, userPost.FirstName, userPost.Surname, userPost.Email, hashedPassword}
	newUser, err = AddUser(newUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Add the new album to the slice.
	c.IndentedJSON(http.StatusCreated, newUser)
}

func getUserById(c *gin.Context) {

	idStr := c.Param("id")

	//pase the number into an unsigned unit64 in base 10
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"message": "User with given Id doesn't exist"})
		return
	}

	c.IndentedJSON(http.StatusOK, user)

}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := GetUserByEmail(creds.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"message": "User with given Id doesn't exist"})
		return
	}

	// Validate credentials here (usually check against a database)
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(creds.Password))
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	InitDB()
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
