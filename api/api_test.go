/*
Package api containing API initialization and API route handler
*/
package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/reyhanfikridz/ecom-account-service/internal/config"
	"github.com/reyhanfikridz/ecom-account-service/internal/utils"
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

// TestInitDB test InitDB
func TestInitDB(t *testing.T) {
	a := API{}

	DBConfig := map[string]string{
		"user":     config.DBUsername,
		"password": config.DBPassword,
		"dbname":   config.DBName,
	}
	err := a.InitDB(DBConfig)
	if err != nil {
		t.Errorf("Expected database connection success,"+
			" but connection failed => %s", err.Error())
	}
}

// TestInitRouter test InitRouter
func TestInitRouter(t *testing.T) {
	a := API{}
	err := a.InitRouter()
	if err != nil {
		t.Errorf("Expected router initialization success,"+
			" but it is failed => %s", err.Error())
	}
}

// TestRegisterHandler test RegisterHandler
func TestRegisterHandler(t *testing.T) {
	// initialize testing API
	a, err := GetTestingAPI()
	if err != nil {
		t.Errorf("There's an error when getting testing API => " + err.Error())
	}

	// initialize testing table
	testTable := []struct {
		FormData        map[string]io.Reader
		ExpectedStatus  int
		ExpectedBodyKey []string
		DeleteDataFirst bool
	}{
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"message", "id"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: false,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader(""),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader(""),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader(""),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader(""),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader(""),
				"role":         strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
		{
			FormData: map[string]io.Reader{
				"email":        strings.NewReader("testregister@gmail.com"),
				"password":     strings.NewReader("test"),
				"full_name":    strings.NewReader("test"),
				"address":      strings.NewReader("test"),
				"phone_number": strings.NewReader("test"),
				"role":         strings.NewReader(""),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
			DeleteDataFirst: true,
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// delete testing data first
		if test.DeleteDataFirst {
			_, err := a.DB.Exec(`DELETE FROM account_user WHERE email = $1`, "testregister@gmail.com")
			if err != nil {
				t.Errorf("There's an error when deleting previous testing data " + err.Error())
			}
		}

		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}

			_, err = io.Copy(fw, r)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}
		}
		w.Close()

		// create new request
		req, err := http.NewRequest("POST", "/api/register/", &bFormData)
		if err != nil {
			t.Errorf("There's an error when creating request API register => " +
				err.Error())
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		// run request
		response := httptest.NewRecorder()
		a.Router.ServeHTTP(response, req)

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("Expected status %d got %d", test.ExpectedStatus, response.Code)
		}

		var responseData map[string]any
		err = json.Unmarshal(response.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("There's an error when unmarshal body response => " + err.Error())
		}
		for _, expectedKey := range test.ExpectedBodyKey {
			if responseData[expectedKey] == nil {
				t.Errorf("Expected key " + expectedKey + " empty/not found")
			}
		}
	}
}

// TestLoginHandler test LoginHandler
func TestLoginHandler(t *testing.T) {
	// initialize testing API
	a, err := GetTestingAPI()
	if err != nil {
		t.Errorf("There's an error when getting testing API => " + err.Error())
	}

	// create user first
	_, err = a.DB.Exec(`DELETE FROM account_user WHERE email = $1`, "testlogin@gmail.com")
	if err != nil {
		t.Errorf("There's an error when deleting login testing data " + err.Error())
	}

	hashedPassword, err := utils.HashPassword("test")
	if err != nil {
		t.Errorf("There's an error when hashing password => " + err.Error())
	}

	createdRow := a.DB.QueryRow(`
		INSERT INTO account_user(email, password, full_name, address, phone_number, role)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		"testlogin@gmail.com", hashedPassword, "test", "test", "test", "test")
	if createdRow.Err() != nil {
		t.Errorf(("There's an error when creating testing user data => " +
			createdRow.Err().Error()))
	}

	// initialize testing table
	testTable := []struct {
		FormData        map[string]io.Reader
		ExpectedStatus  int
		ExpectedBodyKey []string
	}{
		{
			FormData: map[string]io.Reader{
				"email":    strings.NewReader("testlogin@gmail.com"),
				"password": strings.NewReader("test"),
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"message", "token", "role"},
		},
		{
			FormData: map[string]io.Reader{
				"email":    strings.NewReader("testlogin@gmail.co"),
				"password": strings.NewReader("test"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
		},
		{
			FormData: map[string]io.Reader{
				"email":    strings.NewReader("testlogin@gmail.com"),
				"password": strings.NewReader("tes"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}

			_, err = io.Copy(fw, r)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}
		}
		w.Close()

		// create new request
		req, err := http.NewRequest("POST", "/api/login/", &bFormData)
		if err != nil {
			t.Errorf("There's an error when creating request API login => " +
				err.Error())
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		// run request
		response := httptest.NewRecorder()
		a.Router.ServeHTTP(response, req)

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("Expected status %d got %d", test.ExpectedStatus, response.Code)
		}

		var responseData map[string]any
		err = json.Unmarshal(response.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("There's an error when unmarshal body response => " + err.Error())
		}
		for _, expectedKey := range test.ExpectedBodyKey {
			if responseData[expectedKey] == nil {
				t.Errorf("Expected key " + expectedKey + " empty/not found")
			}
		}
	}
}

// TestAuthorizeHandler test AuthorizeHandler
func TestAuthorizeHandler(t *testing.T) {
	// initialize testing API
	a, err := GetTestingAPI()
	if err != nil {
		t.Errorf("There's an error when getting testing API => " + err.Error())
	}

	// create user
	_, err = a.DB.Exec(`DELETE FROM account_user WHERE email = $1`, "testauthorize@gmail.com")
	if err != nil {
		t.Errorf("There's an error when deleting authorize testing data " + err.Error())
	}

	hashedPassword, err := utils.HashPassword("test")
	if err != nil {
		t.Errorf("There's an error when hashing password => " + err.Error())
	}

	createdRow := a.DB.QueryRow(`
		INSERT INTO account_user(email, password, full_name, address, phone_number, role)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		"testauthorize@gmail.com", hashedPassword, "test", "test", "test", "test")
	if createdRow.Err() != nil {
		t.Errorf("There's an error when creating testing user data => " +
			createdRow.Err().Error())
	}

	var userID int
	err = createdRow.Scan(&userID)
	if err != nil {
		t.Errorf("There's an error when creating testing user data => " +
			err.Error())
	}

	// create user session
	validToken, err := utils.GenerateJWT("testauthorize@gmail.com", "test")
	if err != nil {
		t.Errorf("There's an error when creating token " +
			"for creating user session data => " +
			err.Error())
	}

	createdRow = a.DB.QueryRow(`
		INSERT INTO account_usersession(token, account_user_id) 
			VALUES($1,$2) RETURNING id`, validToken, userID)

	if createdRow.Err() != nil {
		t.Errorf("There's an error when creating testing user session data => " +
			createdRow.Err().Error())
	}

	// initialize testing table
	testTable := []struct {
		FormData        map[string]io.Reader
		ExpectedStatus  int
		ExpectedBodyKey []string
	}{
		{
			FormData: map[string]io.Reader{
				"token": strings.NewReader(validToken),
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"id", "email", "password", "address", "phone_number", "role"},
		},
		{
			FormData: map[string]io.Reader{
				"token": strings.NewReader("Invalid Token"),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
		},
		{
			FormData: map[string]io.Reader{
				"token": strings.NewReader(""),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}

			_, err = io.Copy(fw, r)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}
		}
		w.Close()

		// create new request
		req, err := http.NewRequest("POST", "/api/authorize/", &bFormData)
		if err != nil {
			t.Errorf("There's an error when creating request API authorize => " +
				err.Error())
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		// run request
		response := httptest.NewRecorder()
		a.Router.ServeHTTP(response, req)

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("Expected status %d got %d", test.ExpectedStatus, response.Code)
		}

		var responseData map[string]any
		err = json.Unmarshal(response.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("There's an error when unmarshal body response => " + err.Error())
		}
		for _, expectedKey := range test.ExpectedBodyKey {
			if responseData[expectedKey] == nil {
				t.Errorf("Expected key " + expectedKey + " empty/not found")
			}
		}
	}
}

// TestLogoutHandler test LogoutHandler
func TestLogoutHandler(t *testing.T) {
	// initialize testing API
	a, err := GetTestingAPI()
	if err != nil {
		t.Errorf("There's an error when getting testing API => " + err.Error())
	}

	// create user
	_, err = a.DB.Exec(`DELETE FROM account_user WHERE email = $1`, "testlogout@gmail.com")
	if err != nil {
		t.Errorf("There's an error when deleting authorize testing data " + err.Error())
	}

	hashedPassword, err := utils.HashPassword("test")
	if err != nil {
		t.Errorf("There's an error when hashing password => " + err.Error())
	}

	createdRow := a.DB.QueryRow(`
		INSERT INTO account_user(email, password, full_name, address, phone_number, role)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		"testlogout@gmail.com", hashedPassword, "test", "test", "test", "test")
	if createdRow.Err() != nil {
		t.Errorf("There's an error when creating testing user data => " +
			createdRow.Err().Error())
	}

	var userID int
	err = createdRow.Scan(&userID)
	if err != nil {
		t.Errorf("There's an error when creating testing user data => " +
			err.Error())
	}

	// create user session
	validToken, err := utils.GenerateJWT("testlogout@gmail.com", "test")
	if err != nil {
		t.Errorf("There's an error when creating token " +
			"for creating user session data => " +
			err.Error())
	}

	createdRow = a.DB.QueryRow(`
		INSERT INTO account_usersession(token, account_user_id) 
			VALUES($1,$2) RETURNING id`, validToken, userID)

	if createdRow.Err() != nil {
		t.Errorf("There's an error when creating testing user session data => " +
			createdRow.Err().Error())
	}

	// initialize testing table
	testTable := []struct {
		FormData        map[string]io.Reader
		ExpectedStatus  int
		ExpectedBodyKey []string
	}{
		{
			FormData: map[string]io.Reader{
				"token": strings.NewReader(validToken),
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"message"},
		},
		{
			FormData: map[string]io.Reader{
				"token": strings.NewReader(""),
			},
			ExpectedStatus:  400,
			ExpectedBodyKey: []string{"message"},
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// transform form data to bytes buffer
		var bFormData bytes.Buffer
		w := multipart.NewWriter(&bFormData)
		for key, r := range test.FormData {
			fw, err := w.CreateFormField(key)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}

			_, err = io.Copy(fw, r)
			if err != nil {
				t.Errorf("There's an error when creating bytes buffer form data => " +
					err.Error())
			}
		}
		w.Close()

		// create new request
		req, err := http.NewRequest("POST", "/api/logout/", &bFormData)
		if err != nil {
			t.Errorf("There's an error when creating request API logout => " +
				err.Error())
		}
		req.Header.Set("Content-Type", w.FormDataContentType())

		// run request
		response := httptest.NewRecorder()
		a.Router.ServeHTTP(response, req)

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("Expected status %d got %d", test.ExpectedStatus, response.Code)
		}

		var responseData map[string]any
		err = json.Unmarshal(response.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("There's an error when unmarshal body response => " + err.Error())
		}
		for _, expectedKey := range test.ExpectedBodyKey {
			if responseData[expectedKey] == nil {
				t.Errorf("Expected key " + expectedKey + " empty/not found")
			}
		}
	}
}

// TestGetUserHandler test GetUserHandler
func TestGetUserHandler(t *testing.T) {
	// initialize testing API
	a, err := GetTestingAPI()
	if err != nil {
		t.Errorf("There's an error when getting testing API => " + err.Error())
	}

	// create user
	_, err = a.DB.Exec(`DELETE FROM account_user WHERE email = $1`, "test@gmail.com")
	if err != nil {
		t.Errorf("There's an error when deleting authorize testing data " + err.Error())
	}

	hashedPassword, err := utils.HashPassword("test")
	if err != nil {
		t.Errorf("There's an error when hashing password => " + err.Error())
	}

	createdRow := a.DB.QueryRow(`
		INSERT INTO account_user(email, password, full_name, address, phone_number, role)
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		"test@gmail.com", hashedPassword, "test", "test", "test", "test")
	if createdRow.Err() != nil {
		t.Errorf("There's an error when creating testing user data => " +
			createdRow.Err().Error())
	}

	var userID int
	err = createdRow.Scan(&userID)
	if err != nil {
		t.Errorf("There's an error when creating testing user data => " +
			err.Error())
	}

	// initialize testing table
	testTable := []struct {
		Filter          map[string]string
		ExpectedStatus  int
		ExpectedBodyKey []string
	}{
		{
			Filter: map[string]string{
				"email": "test@gmail.com",
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"id", "email", "password", "address", "phone_number", "role"},
		},
		{
			Filter: map[string]string{
				"id": strconv.Itoa(userID),
			},
			ExpectedStatus:  200,
			ExpectedBodyKey: []string{"id", "email", "password", "address", "phone_number", "role"},
		},
	}

	// loop test in test table
	for _, test := range testTable {
		// get url params
		params := url.Values{}
		for key, value := range test.Filter {
			params.Add(key, value)
		}

		// create new request
		req, err := http.NewRequest("GET", "/api/user/", nil)
		req.URL.RawQuery = params.Encode()
		if err != nil {
			t.Errorf("There's an error when creating request API get user => " +
				err.Error())
		}

		// run request
		response := httptest.NewRecorder()
		a.Router.ServeHTTP(response, req)

		// check response
		if response.Code != test.ExpectedStatus {
			t.Errorf("Expected status %d got %d", test.ExpectedStatus, response.Code)
		}

		var responseData map[string]any
		err = json.Unmarshal(response.Body.Bytes(), &responseData)
		if err != nil {
			t.Errorf("There's an error when unmarshal body response => " + err.Error())
		}
		for _, expectedKey := range test.ExpectedBodyKey {
			if responseData[expectedKey] == nil {
				t.Errorf("Expected key " + expectedKey + " empty/not found")
			}
		}
	}

}

// GetTestingAPI get API for testing
func GetTestingAPI() (API, error) {
	a := API{}

	// init database
	DBConfig := map[string]string{
		"user":     config.DBUsername,
		"password": config.DBPassword,
		"dbname":   config.DBTestName,
	}
	err := a.InitDB(DBConfig)
	if err != nil {
		return a, err
	}

	// init router
	err = a.InitRouter()
	if err != nil {
		return a, err
	}

	return a, nil
}
