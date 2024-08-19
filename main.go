package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	logger := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
	// Create a reverse proxy to the users service
	usersServiceURL, err := url.Parse("http://users-service:8080")
	if err != nil {
		log.Fatal(err)
	}
	usersServiceProxy := httputil.NewSingleHostReverseProxy(usersServiceURL)

	// Route to forward requests to the users service
	router.Any("/users/*path", func(c *gin.Context) {
		// Remove the /users prefix before forwarding to the users service
		c.Request.URL.Path = c.Param("path")
		logger.Println(c.Request.URL.Path)
		usersServiceProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Api Gateway"})
	})

	router.Run(":8080")
}
