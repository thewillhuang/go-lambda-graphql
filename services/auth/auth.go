package auth

import (
	"fmt"
	"go-lambda-graphql/config"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var hashCost int

func init() {
	hashCost = 1
	if config.IsProduction {
		hashCost = bcrypt.DefaultCost
	}
}

// HashPassword takes a password and converts it into a one time hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

// CheckPasswordHash takes a password and returns true if it matches with has
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetToken takes a jwt and returns a token struct
func GetToken(Jwt string) (*jwt.Token, error) {
	token, err := jwt.Parse(Jwt, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JWTSecret), nil
	})
	return token, err
}
