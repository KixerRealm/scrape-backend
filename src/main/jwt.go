package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT signing key
var jwtKey = []byte("eyJhbGciOiJIUzI1NiJ9.eyJSb2xlIjoiQWRtaW4iLCJJc3N1ZXIiOiJJc3N1ZXIiLCJVc2VybmFtZSI6IkphdmFJblVzZSIsImV4cCI6MTcwMTUzNTU1MCwiaWF0IjoxNzAxNTM1NTUwfQ.tCMI72EsqfhGjRaOVXMP8V0Bud-gdAzhg-SfDQ2SBG4")

// Claims represents the JWT claims
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateJWT(email string, expirationTime time.Time) (string, error) {
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
