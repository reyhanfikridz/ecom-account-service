/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import (
	"testing"
)

// TestHashPasswordAndComparePassword integration test
// HashPassword and ComparePassword
func TestHashPasswordAndComparePassword(t *testing.T) {
	password := "password"
	failedPassword := "passwprd"

	// hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("There's an error when hashing password => " + err.Error())
	}

	// compare password
	successErr := ComparePassword(hashedPassword, password)
	failedErr := ComparePassword(hashedPassword, failedPassword)

	// check result
	if successErr != nil {
		t.Errorf("Expected password valid, returned password invalid")
	}

	if failedErr == nil {
		t.Errorf("Expected password invalid, returned password valid")
	}
}
