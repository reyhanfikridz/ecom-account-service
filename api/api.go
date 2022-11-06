/*
Package api containing API initialization and API route handler
*/
package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/reyhanfikridz/ecom-account-service/internal/config"
	"github.com/reyhanfikridz/ecom-account-service/internal/form"
	"github.com/reyhanfikridz/ecom-account-service/internal/model"
	"github.com/reyhanfikridz/ecom-account-service/internal/utils"
)

// API contain database connection and router for account service API
type API struct {
	DB     *sql.DB
	Router *mux.Router
}

// InitDB initialize API database connection
func (a *API) InitDB(DBConfig map[string]string) error {
	// connect to db
	connString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DBConfig["user"], DBConfig["password"], DBConfig["dbname"])

	var err error
	a.DB, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	// create db tables if not exist
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

	_, err = a.DB.Exec(tableCreationQuery)
	if err != nil {
		return err
	}

	return nil
}

// InitRouter initialize router for API
func (a *API) InitRouter() error {
	a.Router = mux.NewRouter()

	// route register user
	registerRoute := a.Router.
		HandleFunc("/api/register/", a.RegisterHandler).
		Methods("POST")
	if registerRoute.GetError() != nil {
		return registerRoute.GetError()
	}

	// route login user
	loginRoute := a.Router.
		HandleFunc("/api/login/", a.LoginHandler).
		Methods("POST")
	if loginRoute.GetError() != nil {
		return loginRoute.GetError()
	}

	// route authorize user
	authorizeRoute := a.Router.
		HandleFunc("/api/authorize/", a.AuthorizeHandler).
		Methods("POST")
	if authorizeRoute.GetError() != nil {
		return authorizeRoute.GetError()
	}

	// route logout user
	logoutRoute := a.Router.
		HandleFunc("/api/logout/", a.LogoutHandler).
		Methods("POST")
	if logoutRoute.GetError() != nil {
		return logoutRoute.GetError()
	}

	// route get user
	getUserRoute := a.Router.
		HandleFunc("/api/user/", a.GetUserHandler).
		Methods("GET")
	if getUserRoute.GetError() != nil {
		return getUserRoute.GetError()
	}

	return nil
}

// RegisterHandler handling route register User (method: POST)
func (a *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var responseContent map[string]any
	var responseStatus int

	// allow host
	w.Header().Set("Access-Control-Allow-Origin", config.FrontendURL)

	// get user data from form-data
	u := model.User{
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		FullName:    r.FormValue("full_name"),
		Address:     r.FormValue("address"),
		PhoneNumber: r.FormValue("phone_number"),
		Role:        r.FormValue("role"),
	}

	// validate register user form
	isValid, errString := form.IsUserFormValid(u, "register")
	if isValid { // if register form valid, create user
		u, err := model.CreateUser(a.DB, u)
		if err == nil { // if there's no error when create user
			log.Println(strconv.Quote("POST /api/register/"), "200 SUCCESS")
			responseContent = map[string]any{
				"message": "User registered!",
				"id":      u.ID,
			}
			responseStatus = 200

		} else { // if there's an error when create user
			if strings.Contains(err.Error(), "duplicate") &&
				strings.Contains(err.Error(), "email") { // if email already used
				log.Println(strconv.Quote("POST /api/register/"), "400 BAD REQUEST")
				responseContent = map[string]any{
					"message": "Email already registered, please use another email",
				}
				responseStatus = 400
			} else { // if another error
				log.Println(strconv.Quote("POST /api/register/"), "500 INTERNAL SERVER ERROR")
				log.Println(err.Error())
				responseContent = map[string]any{
					"message": "There's an error when register user => " + err.Error(),
				}
				responseStatus = 500
			}
		}
	} else { // if register form not valid, return bad request
		log.Println(strconv.Quote("POST /api/register/"), "400 BAD REQUEST")
		responseContent = map[string]any{
			"message": errString,
		}
		responseStatus = 400
	}

	response, marshalErr := json.Marshal(responseContent)
	if marshalErr != nil {
		log.Fatal("Error When Creating Response -> ", marshalErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}

// LoginHandler handling route login User (method: POST)
func (a *API) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var responseContent map[string]any
	var responseStatus int

	// allow host
	w.Header().Set("Access-Control-Allow-Origin", config.FrontendURL)

	// get user data from form-data
	u := model.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// validate user data from login form
	isValid, errString := form.IsUserFormValid(u, "login")
	if isValid { // if user data valid, login user
		token, status, u, err := model.AuthenticateUser(a.DB, u)
		if status == 200 && err == nil { // If user authenticated
			log.Println(strconv.Quote("POST /api/login/"), "200 SUCCESS")
			responseContent = map[string]any{
				"message": "User logged in!",
				"token":   token,
				"role":    u.Role,
			}
			responseStatus = 200
		} else if status == 400 { // if user not authenticated
			log.Println(strconv.Quote("POST /api/login/"), "400 BAD REQUEST")
			responseContent = map[string]any{
				"message": "Email or Password invalid",
			}
			responseStatus = 400
		} else { // if there's internal server error
			log.Println(strconv.Quote("POST /api/login/"), "500 INTERNAL SERVER ERROR")
			log.Println(err.Error())
			responseContent = map[string]any{
				"message": err,
			}
			responseStatus = 500
		}

	} else { // if user data from login form not valid, return bad request
		log.Println(strconv.Quote("POST /api/login/"), "400 BAD REQUEST")
		responseContent = map[string]any{
			"message": errString,
		}
		responseStatus = 400
	}

	response, marshalErr := json.Marshal(responseContent)
	if marshalErr != nil {
		log.Fatal("Error When Creating Response -> ", marshalErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}

// AuthorizeHandler handling route authorize user (method: POST)
func (a *API) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	var responseStatus int
	var isResponseData bool
	responseMessage := make(map[string]string)
	responseData := model.User{}

	// allow host
	w.Header().Set("Access-Control-Allow-Origin", config.FrontendURL)

	// check token in form or not
	tokenString := r.FormValue("token")
	if strings.TrimSpace(tokenString) != "" { // if token is in form

		// validate token
		tokenClaimsMap := utils.ValidateJWT(tokenString)
		if tokenClaimsMap != nil { // if token valid

			// check if user is in DB
			user, err := model.GetUser(a.DB, tokenClaimsMap["email"], 0)
			if err == nil { // if user exist

				// check if user session is in DB
				userSession, err := model.GetUserSession(a.DB, tokenString, user.ID)
				if err == nil && userSession.ID != 0 { // if user session exist
					log.Println(strconv.Quote("POST /api/authorize/"), "200 SUCCESS")

					user.Password = "" // makes password empty for security purpose

					isResponseData = true
					responseData = user
					responseStatus = 200
				} else if err == nil && userSession.ID == 0 { // if user session not exist
					log.Println(strconv.Quote("POST /api/authorize/"), "400 BAD REQUEST")

					isResponseData = false
					responseMessage["message"] = "Token not valid"
					responseStatus = 400
				} else { // if error encountered
					log.Println(strconv.Quote("POST /api/authorize/"), "500 INTERNAL SERVER ERROR")
					log.Println(err.Error())

					isResponseData = false
					responseMessage["message"] = err.Error()
					responseStatus = 500
				}

			} else { // if user not exist
				log.Println(strconv.Quote("POST /api/authorize/"), "500 INTERNAL SERVER ERROR")
				log.Println(err.Error())

				isResponseData = false
				responseMessage["message"] = err.Error()
				responseStatus = 500
			}

		} else { // if token not valid
			log.Println(strconv.Quote("POST /api/authorize/"), "400 BAD REQUEST")

			isResponseData = false
			responseMessage["message"] = "Token not valid"
			responseStatus = 400
		}

	} else { // if token not in form
		log.Println(strconv.Quote("POST /api/authorize/"), "400 BAD REQUEST")

		isResponseData = false
		responseMessage["message"] = "Token empty/not found"
		responseStatus = 400
	}

	// write response
	var response []byte
	var marshalErr error
	if isResponseData {
		response, marshalErr = json.Marshal(responseData)
	} else {
		response, marshalErr = json.Marshal(responseMessage)
	}

	if marshalErr != nil {
		log.Fatal("Error When Creating Response -> ", marshalErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}

// LogoutHandler handling route logout user (method: POST)
func (a *API) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var responseContent map[string]any
	var responseStatus int

	// allow host
	w.Header().Set("Access-Control-Allow-Origin", config.FrontendURL)

	// check token in form
	tokenString := r.FormValue("token")
	if strings.TrimSpace(r.FormValue("token")) != "" { // if token exist
		// delete user session
		err := model.DeleteUserSession(a.DB, tokenString)
		if err == nil { // if delete success
			log.Println(strconv.Quote("POST /api/logout/"), "200 SUCCESS")
			responseContent = map[string]any{
				"message": "User logged out",
			}
			responseStatus = 200
		} else { // if there's an error
			log.Println(strconv.Quote("POST /api/logout/"), "500 INTERNAL SERVER ERROR")
			log.Println(err.Error())
			responseContent = map[string]any{
				"message": err,
			}
			responseStatus = 500
		}
	} else { // if token not exist
		log.Println(strconv.Quote("POST /api/logout/"), "400 BAD REQUEST")
		responseContent = map[string]any{
			"message": "Token empty/not found",
		}
		responseStatus = 400
	}

	response, marshalErr := json.Marshal(responseContent)
	if marshalErr != nil {
		log.Fatal("Error When Creating Response -> ", marshalErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}

// GetUserHandler handling route get user data (method: GET)
func (a *API) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var responseStatus int
	var isResponseData bool
	responseMessage := make(map[string]string)
	responseData := model.User{}

	// allow host
	w.Header().Set("Access-Control-Allow-Origin", config.ProductServiceURL)

	// check id in params or not
	stringID := r.FormValue("id")
	if strings.TrimSpace(stringID) != "" { // if token is in form

		// convert id to integer
		ID, err := strconv.Atoi(stringID)
		if err == nil { // if convert success

			// get user
			user, err := model.GetUser(a.DB, "", ID)
			if err == nil { // if get user success
				log.Println(strconv.Quote("GET /api/user/"), "200 SUCCESS")

				user.Password = "" // makes password empty for security purpose

				isResponseData = true
				responseData = user
				responseStatus = 200

			} else { // if get user failed
				log.Println(strconv.Quote("GET /api/user/"), "500 INTERNAL SERVER ERROR")
				log.Println(err.Error())

				isResponseData = false
				responseMessage["message"] = err.Error()
				responseStatus = 500
			}

		} else { // if ID not valid
			log.Println(strconv.Quote("GET /api/user/"), "400 BAD REQUEST")
			log.Println(err.Error())

			isResponseData = false
			responseMessage["message"] = err.Error()
			responseStatus = 400
		}

	} else { // if id not in params

		// check email in params
		email := r.FormValue("email")
		if strings.TrimSpace(email) != "" { // if email valid

			// get user
			user, err := model.GetUser(a.DB, email, 0)
			if err == nil { // if get user success
				log.Println(strconv.Quote("GET /api/user/"), "200 SUCCESS")

				user.Password = "" // makes password empty for security purpose

				isResponseData = true
				responseData = user
				responseStatus = 200

			} else { // if get user failed
				log.Println(strconv.Quote("GET /api/user/"), "500 INTERNAL SERVER ERROR")
				log.Println(err.Error())

				isResponseData = false
				responseMessage["message"] = err.Error()
				responseStatus = 500
			}

		} else { // if email not valid
			log.Println(strconv.Quote("GET /api/user/"), "400 BAD REQUEST")

			isResponseData = false
			responseMessage["message"] = "email empty/not found"
			responseStatus = 400
		}

	}

	// write response
	var response []byte
	var marshalErr error
	if isResponseData {
		response, marshalErr = json.Marshal(responseData)
	} else {
		response, marshalErr = json.Marshal(responseMessage)
	}

	if marshalErr != nil {
		log.Fatal("Error When Creating Response -> ", marshalErr)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(response)
}
