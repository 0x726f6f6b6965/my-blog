package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generates a new JWT token
func GenerateJWT(email string, secret string, expire time.Duration) (string, error) {
	// create a JWT claim
	claims := jwt.MapClaims{}
	// assign an expiration time for the token
	claims["exp"] = time.Now().Add(expire).Unix()
	// assign a data for email
	claims["user"] = email

	// create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// convert the JWT token into the string
	t, err := token.SignedString([]byte(secret))

	// if conversion is failed, return an error
	if err != nil {
		return "", err
	}

	// return the generated JWT token
	return t, nil
}
