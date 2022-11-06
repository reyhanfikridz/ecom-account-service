/*
Package utils containing utilities function

This package cannot have import from another package except for config package
*/
package utils

import (
	"log"
	"testing"

	"github.com/reyhanfikridz/ecom-account-service/internal/config"
)

// TestMain do some test before and after all testing in the package
func TestMain(m *testing.M) {
	// init all config before can be used
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("There's an error when initialize config => %s", err)
	}

	// run all testing
	m.Run()
}

// TestGenerateJWTAndValidateJWT integration test
// GenerateJWT and ValidateJWT
func TestGenerateJWTAndValidateJWT(t *testing.T) {
	email := "admin@gmail.com"
	role := "admin"

	// generate jwt
	tokenString, err := GenerateJWT(email, role)
	if err != nil {
		t.Errorf("There's an error when generate JWT => " + err.Error())
	}

	// validate jwt
	tokenClaimsMap := ValidateJWT(tokenString)

	// check result
	if tokenClaimsMap == nil {
		t.Errorf("Expected JWT token valid, but got invalid")
	}

	if email != tokenClaimsMap["email"] {
		t.Errorf("Expected email '" + email + "', but got '" + tokenClaimsMap["email"] + "'")
	}

	if role != tokenClaimsMap["role"] {
		t.Errorf("Expected role '" + role + "', but got '" + tokenClaimsMap["role"] + "'")
	}
}
