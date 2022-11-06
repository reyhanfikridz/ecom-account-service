/*
Package main the executeable file
*/
package main

import (
	"log"
	"net/http"

	"github.com/reyhanfikridz/ecom-account-service/api"
	"github.com/reyhanfikridz/ecom-account-service/internal/config"
)

// main
func main() {
	// init API
	a, err := InitAPI()
	if err != nil {
		log.Fatal(err)
	}

	// serve server
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

// InitAPI initialize API
func InitAPI() (api.API, error) {
	a := api.API{}

	// init all config before can be used
	err := config.InitConfig()
	if err != nil {
		return a, err
	}

	// init database
	DBConfig := map[string]string{
		"user":     config.DBUsername,
		"password": config.DBPassword,
		"dbname":   config.DBName,
	}
	err = a.InitDB(DBConfig)
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
