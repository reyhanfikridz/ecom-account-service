/*
Package form collection of form validation
*/
package form

import (
	"strings"

	"github.com/reyhanfikridz/ecom-account-service/internal/model"
)

// IsUserFormValid check if user form is valid (for login/register)
func IsUserFormValid(u model.User, formType string) (bool, string) {
	if strings.TrimSpace(u.Email) == "" {
		return false, "email empty/not found"
	}

	if strings.TrimSpace(u.Password) == "" {
		return false, "password empty/not found"
	}

	if formType == "register" && strings.TrimSpace(u.FullName) == "" {
		return false, "full_name empty/not found"
	}

	if formType == "register" && strings.TrimSpace(u.Address) == "" {
		return false, "address empty/not found"
	}

	if formType == "register" && strings.TrimSpace(u.PhoneNumber) == "" {
		return false, "phone_number empty/not found"
	}

	if formType == "register" && strings.TrimSpace(u.Role) == "" {
		return false, "role empty/not found"
	}

	return true, ""
}
