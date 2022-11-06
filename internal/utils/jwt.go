/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import (
	"fmt"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/reyhanfikridz/ecom-account-service/internal/config"
)

// GenerateJWT generate jwt token string
func GenerateJWT(email string, role string) (string, error) {
	// initialize new token
	token := jwt.New(config.JWTSigningMethod)

	// set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// get token string
	tokenString, err := token.SignedString([]byte(config.JWTSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validate jwt token string
//
// If token not valid, return nil
func ValidateJWT(tokenString string) map[string]string {
	// parse token from token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token not using the right signing method")
		}
		return []byte(config.JWTSecretKey), nil
	})
	if err != nil {
		return nil
	}

	// Check and get token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}

	// Check if token not valid
	if !token.Valid {
		return nil
	}

	// check claims data
	claimsMap := map[string]string{
		"email": fmt.Sprint(claims["email"]),
		"role":  fmt.Sprint(claims["role"]),
	}
	for _, item := range claimsMap {
		if strings.TrimSpace(item) == "" {
			return nil
		}
	}

	return claimsMap
}
