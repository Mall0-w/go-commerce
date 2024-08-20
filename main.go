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
		// Log the incoming request for debugging
		logger.Printf("Forwarding request: %s %s\n", c.Request.Method, c.Request.URL.String())

		// Forward the request as-is to the users service
		usersServiceProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Api Gateway"})
	})

	router.Run(":8080")
}
