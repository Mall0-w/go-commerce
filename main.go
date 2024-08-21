package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

var logger *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

func getProxy(proxyUrl string) func(*gin.Context) {
	serviceUrl, err := url.Parse(proxyUrl)
	if err != nil {
		logger.Fatal(err)
	}
	//define proxies
	serviceProxy := httputil.NewSingleHostReverseProxy(serviceUrl)

	return func(c *gin.Context) {
		// Log the incoming request for debugging
		logger.Printf("Forwarding request: %s %s\n", c.Request.Method, c.Request.URL.String())

		// Forward the request as-is to the users service
		serviceProxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()

	// Route to forward requests to the users service
	router.Any("/users/*path", getProxy("http://users-service:8080"))

	router.Any("/products/*path", getProxy("http://product-service:8080"))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Api Gateway"})
	})

	router.Run(":8080")
}
