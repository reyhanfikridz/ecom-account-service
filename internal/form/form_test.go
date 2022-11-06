/*
Package form collection of form validation
*/
package form

import (
	"testing"

	"github.com/reyhanfikridz/ecom-account-service/internal/model"
)

// TestIsUserFormValid test IsUserFormValid
func TestIsUserFormValid(t *testing.T) {
	// initialize testing table
	testTable := []struct {
		Name              string
		User              model.User
		FormType          string
		ExpectedIsValid   bool
		ExpectedErrString string
	}{
		{
			Name: "test-register-success",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "test",
				Address:     "test",
				PhoneNumber: "test",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   true,
			ExpectedErrString: "",
		},
		{
			Name: "test-register-failed-1",
			User: model.User{
				Email:       "",
				Password:    "test",
				FullName:    "test",
				Address:     "test",
				PhoneNumber: "test",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "email empty/not found",
		},
		{
			Name: "test-register-failed-2",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "",
				FullName:    "test",
				Address:     "test",
				PhoneNumber: "test",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "password empty/not found",
		},
		{
			Name: "test-register-failed-3",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "",
				Address:     "test",
				PhoneNumber: "test",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "full_name empty/not found",
		},
		{
			Name: "test-register-failed-4",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "test",
				Address:     "",
				PhoneNumber: "test",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "address empty/not found",
		},
		{
			Name: "test-register-failed-5",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "test",
				Address:     "test",
				PhoneNumber: "",
				Role:        "test",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "phone_number empty/not found",
		},
		{
			Name: "test-register-failed-6",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "test",
				Address:     "test",
				PhoneNumber: "test",
				Role:        "",
			},
			FormType:          "register",
			ExpectedIsValid:   false,
			ExpectedErrString: "role empty/not found",
		},
		{
			Name: "test-login-success",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "test",
				FullName:    "",
				Address:     "",
				PhoneNumber: "",
				Role:        "",
			},
			FormType:          "login",
			ExpectedIsValid:   true,
			ExpectedErrString: "",
		},
		{
			Name: "test-login-failed-1",
			User: model.User{
				Email:       "",
				Password:    "test",
				FullName:    "",
				Address:     "",
				PhoneNumber: "",
				Role:        "",
			},
			FormType:          "login",
			ExpectedIsValid:   false,
			ExpectedErrString: "email empty/not found",
		},
		{
			Name: "test-login-failed-2",
			User: model.User{
				Email:       "test@gmail.com",
				Password:    "",
				FullName:    "",
				Address:     "",
				PhoneNumber: "",
				Role:        "",
			},
			FormType:          "login",
			ExpectedIsValid:   false,
			ExpectedErrString: "password empty/not found",
		},
	}

	// loop test in test table
	for _, test := range testTable {
		isValid, errString := IsUserFormValid(test.User, test.FormType)
		if test.ExpectedIsValid && !isValid {
			t.Errorf("Expected form Valid got Invalid")
		} else if !test.ExpectedIsValid && isValid {
			t.Errorf("Expected form Invalid got Valid")
		}

		if test.ExpectedErrString != errString {
			t.Errorf("Expected error '" + test.ExpectedErrString + "' got '" + errString + "'")
		}
	}
}
