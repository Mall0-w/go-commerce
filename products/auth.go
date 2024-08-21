package main

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	gin "github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

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
