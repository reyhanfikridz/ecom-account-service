/*
Package model containing structs and functions for
database transaction
*/
package model

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/reyhanfikridz/ecom-account-service/internal/utils"
)

// user model
type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// user session model
type UserSession struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	User  User   `json:"user"`
}

// func for creating exactly one new user
func CreateUser(DB *sql.DB, u User) (User, error) {
	// hashing password
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return u, err
	}

	//////////////////// begin transaction /////////////////////
	tx, err := DB.Begin()
	if err != nil {
		return u, err
	}
	defer tx.Rollback() // rollback transaction if fail (return before commit)

	// create row data
	createdRow := tx.QueryRow(`
		INSERT INTO account_user(email, password, full_name, address, phone_number, role) 
			VALUES($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		u.Email, hashedPassword, u.FullName, u.Address, u.PhoneNumber, u.Role)

	if createdRow.Err() != nil {
		return u, createdRow.Err()
	}

	// scan row ID into user.ID
	err = createdRow.Scan(&u.ID)
	if err != nil {
		return u, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return u, err
	}
	////////////////////////////////////////////////////////////

	return u, nil
}

// func for authenticate user
func AuthenticateUser(DB *sql.DB, u User) (
	string, int, User, error) {
	// get existed user data
	existedUser, err := GetUser(DB, u.Email, 0)
	if err != nil {
		return "", 400, existedUser, nil
	}

	// check password right or wrong
	err = utils.ComparePassword(existedUser.Password, u.Password)
	if err != nil {
		return "", 400, existedUser, nil
	}

	// generate jwt token string
	tokenString, err := utils.GenerateJWT(existedUser.Email, existedUser.Role)
	if err != nil {
		return "", 500, existedUser, err
	}

	// save user sessions
	userSession := UserSession{
		Token: tokenString,
		User:  existedUser,
	}
	userSession, err = CreateUserSession(DB, userSession)
	if err != nil {
		return "", 500, existedUser, err
	}

	return tokenString, 200, existedUser, nil
}

// func for get user data by email or by ID
func GetUser(DB *sql.DB, email string, ID int) (User, error) {
	user := User{}

	// do query
	var queryRes *sql.Row
	if ID == 0 {
		queryRes = DB.QueryRow(`
			SELECT id, email, password, full_name, address, phone_number, role 
				FROM account_user WHERE email = $1`,
			email)
	} else {
		queryRes = DB.QueryRow(`
			SELECT id, email, password, full_name, address, phone_number, role 
				FROM account_user WHERE id = $1`,
			ID)
	}

	if queryRes.Err() != nil {
		return user, queryRes.Err()
	}

	// get data from query result
	// Note: The order of Scan need to be same as order in QueryRow
	err := queryRes.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Address,
		&user.PhoneNumber,
		&user.Role,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// func for create user session
func CreateUserSession(DB *sql.DB, us UserSession) (UserSession, error) {
	//////////////////// begin transaction /////////////////////
	tx, err := DB.Begin()
	if err != nil {
		return us, err
	}
	defer tx.Rollback() // rollback transaction if fail (return before commit)

	// delete existed user session
	_, err = tx.Exec(`
		DELETE FROM account_usersession
			WHERE account_user_id = $1`, us.User.ID)
	if err != nil {
		return us, err
	}

	// insert user session
	createdRow := tx.QueryRow(`
		INSERT INTO account_usersession(token, account_user_id) 
			VALUES($1,$2) RETURNING id`, us.Token, us.User.ID)
	if createdRow.Err() != nil {
		return us, createdRow.Err()
	}

	// scan row ID into user session ID
	err = createdRow.Scan(&us.ID)
	if err != nil {
		return us, err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return us, err
	}
	////////////////////////////////////////////////////////////

	return us, nil
}

// func for get session data by token and user id
func GetUserSession(DB *sql.DB, tokenString string, userID int) (UserSession, error) {
	userSession := UserSession{}

	// do query
	row := DB.QueryRow(`
		SELECT *
			FROM account_usersession FULL OUTER JOIN account_user
				ON account_usersession.account_user_id = account_user.id
			WHERE token = $1 AND account_user_id = $2
		`, tokenString, userID)

	if row.Err() != nil {
		return userSession, row.Err()
	}

	// scan query result
	err := row.Scan(
		&userSession.ID,
		&userSession.Token,
		&userSession.User.ID,
		&userSession.User.ID,
		&userSession.User.Email,
		&userSession.User.Password,
		&userSession.User.FullName,
		&userSession.User.Address,
		&userSession.User.PhoneNumber,
		&userSession.User.Role,
	)
	if err != nil {
		return userSession, err
	}
	return userSession, nil
}

// func for delete user session by token string
func DeleteUserSession(DB *sql.DB, tokenString string) error {
	//////////////////// begin transaction /////////////////////
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // rollback transaction if fail (return before commit)

	// delete user session
	_, err = tx.Exec(`
		DELETE FROM account_usersession
			WHERE token = $1
	`, tokenString)

	if err != nil {
		return err
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	////////////////////////////////////////////////////////////

	return nil
}
