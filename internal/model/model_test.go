/*
Package model containing structs and functions for
database transaction
*/
package model

import (
	"database/sql"
	"fmt"
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

// TestCreateUser test CreateUser
func TestCreateUser(t *testing.T) {
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous testing data => " +
			err.Error())
	}

	// create user data on DB
	user, err = CreateUser(DB, user)

	// check result
	if err != nil {
		t.Errorf("Expected error nil, but got not nil => " + err.Error())
	}

	if user.ID == 0 {
		t.Errorf("Expected returned ID not zero, but got zero")
	}
}

// TestCreateUserAndGetUser integration test
// CreateUser and GetUser
func TestCreateUserAndGetUser(t *testing.T) {
	//////////////////// CREATE ////////////////////
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous testing data => " +
			err.Error())
	}

	// create user data on DB
	user, _ = CreateUser(DB, user)

	//////////////////// GET BY EMAIL ////////////////////
	userByEmail, err := GetUser(DB, "admin@gmail.com", 0)
	if err != nil {
		t.Errorf("Expected err nil, but got err not nil => " + err.Error())
	}
	if user.ID != userByEmail.ID {
		t.Errorf("Expected user ID %d, but got user ID %d", user.ID, userByEmail.ID)
	}

	//////////////////// GET BY ID ////////////////////
	userByID, err := GetUser(DB, "", user.ID)
	if err != nil {
		t.Errorf("Expected err nil, but got err not nil => " + err.Error())
	}
	if user.ID != userByID.ID {
		t.Errorf("Expected user ID %d, but got user ID %d", user.ID, userByEmail.ID)
	}
}

// TestCreateUserAndCreateUserSession integration test
// CreateUser and CreateUserSession
func TestCreateUserAndCreateUserSession(t *testing.T) {
	//////////////////// CREATE USER ////////////////////
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous user testing data => " +
			err.Error())
	}

	// create user data on DB
	user, _ = CreateUser(DB, user)

	//////////////////// CREATE USER SESSION ////////////////////
	// create user session
	userSession := UserSession{
		Token: "This is token",
		User:  user,
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_usersession WHERE account_user_id = $1`, user.ID)
	if err != nil {
		t.Errorf("There's an error when deleting " +
			"previous user session testing data => " +
			err.Error())
	}

	// create user session data on DB
	userSession, err = CreateUserSession(DB, userSession)

	// check result
	if err != nil {
		t.Errorf("Expected error nil, but got not nil => " + err.Error())
	}

	if userSession.ID == 0 {
		t.Errorf("Expected returned ID not zero, but got zero")
	}
}

// TestCreateUserAndCreateUserSessionAndGetUserSession integration test
// CreateUser, CreateUserSession and GetUserSession
func TestCreateUserAndCreateUserSessionAndGetUserSession(t *testing.T) {
	//////////////////// CREATE USER ////////////////////
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous user testing data => " +
			err.Error())
	}

	// create user data on DB
	user, _ = CreateUser(DB, user)

	//////////////////// CREATE USER SESSION ////////////////////
	// create user session
	userSession := UserSession{
		Token: "This is token",
		User:  user,
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_usersession WHERE account_user_id = $1`, user.ID)
	if err != nil {
		t.Errorf("There's an error when deleting " +
			"previous user session testing data => " +
			err.Error())
	}

	// create user session data on DB
	userSession, _ = CreateUserSession(DB, userSession)

	//////////////////// GET USER SESSION ////////////////////
	userSession2, err := GetUserSession(DB, userSession.Token, userSession.User.ID)
	if err != nil {
		t.Errorf("Expected err nil, but got err not nil => " + err.Error())
	}
	if userSession.ID != userSession2.ID {
		t.Errorf("Expected user session ID %d, but got user session ID %d",
			userSession.ID, userSession2.ID)
	}
	if userSession.User.ID != userSession2.User.ID {
		t.Errorf("Expected user ID %d on user session, but got user ID %d",
			userSession.User.ID, userSession2.User.ID)
	}
}

// TestCreateUserAndCreateUserSessionAndDeleteUserSession integration test
// CreateUser, CreateUserSession and DeleteUserSession
func TestCreateUserAndCreateUserSessionAndDeleteUserSession(t *testing.T) {
	//////////////////// CREATE USER ////////////////////
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous user testing data => " +
			err.Error())
	}

	// create user data on DB
	user, _ = CreateUser(DB, user)

	//////////////////// CREATE USER SESSION ////////////////////
	// create user session
	userSession := UserSession{
		Token: "This is token",
		User:  user,
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_usersession WHERE account_user_id = $1`, user.ID)
	if err != nil {
		t.Errorf("There's an error when deleting " +
			"previous user session testing data => " +
			err.Error())
	}

	// create user session data on DB
	userSession, _ = CreateUserSession(DB, userSession)

	//////////////////// DELETE USER SESSION ////////////////////
	err = DeleteUserSession(DB, userSession.Token)
	if err != nil {
		t.Errorf("Expected err nil (delete success), " +
			"but got err not nil (delete failed) => " + err.Error())
	}

	_, err = GetUserSession(DB, userSession.Token, userSession.User.ID)
	if err == nil {
		t.Errorf("Expected err not nil, but got err nil")
	}
}

// TestCreateUserAndAuthenticateUser integration test
// CreateUser and AuthenticateUser
func TestCreateUserAndAuthenticateUser(t *testing.T) {
	//////////////////// CREATE USER ////////////////////
	// create user
	user := User{
		Email:       "admin@gmail.com",
		Password:    "admin",
		FullName:    "admin",
		Address:     "address",
		PhoneNumber: "08111111111",
		Role:        "admin",
	}

	// get connection to testing DB
	DB, err := GetTestDBConnection()
	if err != nil {
		t.Errorf("Connection to testing DB failed => " + err.Error())
	}

	// delete prev data first
	_, err = DB.Exec(`DELETE FROM account_user WHERE email = $1`, user.Email)
	if err != nil {
		t.Errorf("There's an error when deleting previous user testing data => " +
			err.Error())
	}

	// create user data on DB
	user, _ = CreateUser(DB, user)

	//////////////////// AUTHENTICATE USER ////////////////////
	// create testing table
	testTable := []struct {
		User           User
		ExpectedStatus int
	}{
		{
			User: User{
				Email:    "admin@gmail.com",
				Password: "admin",
			},
			ExpectedStatus: 200,
		},
		{
			User: User{
				Email:    "admin@gmail.co",
				Password: "admin",
			},
			ExpectedStatus: 400,
		},
		{
			User: User{
				Email:    "admin@gmail.com",
				Password: "adminn",
			},
			ExpectedStatus: 400,
		},
	}

	for _, test := range testTable {
		_, status, _, err := AuthenticateUser(DB, test.User)
		if err != nil {
			t.Errorf("There's an error when authenticate user =>" + err.Error())
		}

		if test.ExpectedStatus != status {
			t.Errorf("Expected status %d, but got status %d",
				test.ExpectedStatus, status)
		}
	}
}

// GetTestDBConnection get connection to testing DB
func GetTestDBConnection() (*sql.DB, error) {
	// Connect to db
	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DBUsername, config.DBPassword, config.DBTestName)

	DB, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Create user tables if not exist
	tableCreationQuery := `
		CREATE TABLE IF NOT EXISTS account_user
		(
			id SERIAL PRIMARY KEY NOT NULL,
			email VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			full_name VARCHAR(50) NOT NULL,
			address VARCHAR(100) NOT NULL,
			phone_number VARCHAR(20) NOT NULL,
			role VARCHAR(20) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS account_usersession
		(
			id SERIAL PRIMARY KEY NOT NULL,
			token VARCHAR(200) UNIQUE NOT NULL,
			account_user_id INT UNIQUE NOT NULL,
			CONSTRAINT fk_account_user
				FOREIGN KEY(account_user_id) 
					REFERENCES account_user(id)
					ON DELETE CASCADE
		);
	`

	_, err = DB.Exec(tableCreationQuery)
	if err != nil {
		return nil, err
	}

	return DB, nil
}
