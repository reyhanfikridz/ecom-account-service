/*
Package config collection of configuration
*/
package config

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var (
	DBName     string
	DBTestName string
	DBUsername string
	DBPassword string

	JWTSecretKey     string
	JWTSigningMethod *jwt.SigningMethodHMAC

	ServerURL         string
	FrontendURL       string
	ProductServiceURL string
)

// InitConfig initialize all config variable from environment variable
func InitConfig() error {
	// load all values from .env file into the system
	// .env file must be at root directory (same level as go.mod file)
	err := godotenv.Load(os.ExpandEnv(
		"$GOPATH/src/github.com/reyhanfikridz/ecom-account-service/.env"))
	if err != nil {
		return err
	}

	// set all config variable after all environment variable loaded
	DBName = os.Getenv("ECOM_ACCOUNT_SERVICE_DB_NAME")
	DBTestName = os.Getenv("ECOM_ACCOUNT_SERVICE_DB_TEST_NAME")
	DBUsername = os.Getenv("ECOM_ACCOUNT_SERVICE_DB_USERNAME")
	DBPassword = os.Getenv("ECOM_ACCOUNT_SERVICE_DB_PASSWORD")

	JWTSecretKey = os.Getenv("ECOM_ACCOUNT_SERVICE_JWT_SECRET_KEY")
	JWTSigningMethod = jwt.SigningMethodHS256

	ServerURL = os.Getenv("ECOM_ACCOUNT_SERVICE_URL")
	FrontendURL = os.Getenv("ECOM_ACCOUNT_SERVICE_FRONTEND_URL")
	ProductServiceURL = os.Getenv("ECOM_ACCOUNT_SERVICE_PRODUCT_SERVICE_URL")

	return nil
}
