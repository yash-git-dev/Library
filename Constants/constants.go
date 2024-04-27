package constants

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/golang-jwt/jwt"
)

type User struct {
	Username string
	Password string
	IsAdmin  bool
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var JWTSecretKey = []byte("dfhdgfhwue84837r4hmtknshvagtor84572m")

var Users []User

func InitialiseUsers() {
	Users = []User{
		{Username: "demouser1", Password: "password1", IsAdmin: false},
		{Username: "demouser2", Password: "password2", IsAdmin: false},
		{Username: "demoadmin1", Password: "adminpassword1", IsAdmin: true},
		{Username: "demoadmin2", Password: "adminpassword2", IsAdmin: true},
	}

	for i := range Users {
		Users[i].Password = HashPassword(Users[i].Password)
	}
}

func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	return hashedPassword
}
