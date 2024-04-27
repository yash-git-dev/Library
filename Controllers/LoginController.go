package controllers

import (
	"net/http"
	constants "root/Constants"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserLoginPOST struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var user UserLoginPOST
	var token string
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "please check input")
	}

	for _, userDB := range constants.Users {
		password := constants.HashPassword(user.Password)

		if userDB.Username == user.Username && userDB.Password == password {
			token, err = GenerateToken(userDB.Username)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error while generating token:" + err.Error()})
				return
			}

			c.Header("Authorisation", token)

			c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
}

func GenerateToken(userName string) (string, error) {
	claims := constants.CustomClaims{
		Username: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(constants.JWTSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
