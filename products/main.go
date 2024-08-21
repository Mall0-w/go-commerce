package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var Logger *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

func main() {

	InitDB()
	router := gin.Default()

	router.GET("products/testConnection", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Users Microservice"})
	})

	router.Run(":8080")
}
