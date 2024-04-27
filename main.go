package main

import (
	"fmt"
	"log"
	"net/http"
	constants "root/Constants"
	controllers "root/Controllers"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &constants.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return constants.JWTSecretKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*constants.CustomClaims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", fmt.Errorf("invalid token")
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		username, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !userExists(username) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}

func userExists(username string) bool {
	for _, user := range constants.Users {
		if user.Username == username {
			return true
		}
	}
	return false
}

func main() {
	router := gin.Default()

	constants.InitialiseUsers()

	authorisedRouter := router.Group("")

	router.POST("/login", controllers.Login)

	authorisedRouter.Use(JWTMiddleware())

	authorisedRouter.GET("/home", controllers.BooksGET)
	authorisedRouter.POST("/addBook", controllers.BookPOST)
	authorisedRouter.DELETE("/deleteBook/:bookName", controllers.BookDELETE)

	if err := router.Run("0.0.0.0:8000"); err != nil {
		log.Fatal("Failed to start server at port 8000")
	}
}
