/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword hashing password
func HashPassword(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(hash), err
}

// ComparePassword compare hashed password and inputed password
func ComparePassword(hp string, p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(p))
	return err
}
