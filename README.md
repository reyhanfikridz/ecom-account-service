# ecom-account-service

### ECOM summary:
ECOM is a simple E-Commerce website builded with Go backend microservices and Django frontend. Disclaimer! I have zero real experience in building E-Commerce system, so if the system is really bad, I apologized in advance. This is just my personal project using Go microservices. You can use all the code of this project as a template for real E-Commerce in the future if you like it. Disclaimer again! I also not a frontend specialist, so I just use a free template I found in the internet and an original bootstrap template.

### Repository summary:
This is a microservice for ECOM that related to customer account CRUD.

### Requirements:
1. go (recommended: v1.18.4)
2. postgresql (recommended: v13.4)

### Steps to run the server:
1. install all requirements
2. clone repository with `git clone https://github.com/reyhanfikridz/ecom-account-service` at directory `$GOPATH/src/github.com/reyhanfikridz/`
3. change branch to release-1 with `git checkout release-1` then `git pull origin release-1`
4. install required go library with `go mod download` then `go mod vendor` at repository root directory (same level as README.md)
5. create file .env at repository root directory (same level as README.md) with contents:

```
ECOM_ACCOUNT_SERVICE_DB_NAME=<database name, example: ecom_account_service>
ECOM_ACCOUNT_SERVICE_DB_TEST_NAME=<database test name, example: ecom_account_service_test>
ECOM_ACCOUNT_SERVICE_DB_USERNAME=<postgres username>
ECOM_ACCOUNT_SERVICE_DB_PASSWORD=<postgres password>

ECOM_ACCOUNT_SERVICE_JWT_SECRET_KEY=<jwt secret key, example: O3u12Kb0ciRNCOrDjR4E5YuNoDjIR95FtPjFxK0DqVHKgQcadvkK5UnB2OeLeQOa>

ECOM_ACCOUNT_SERVICE_URL=<this service url, example: :8010>
ECOM_ACCOUNT_SERVICE_FRONTEND_URL=<ecom frontend url, example: http://127.0.0.1:8000>
ECOM_ACCOUNT_SERVICE_PRODUCT_SERVICE_URL=<ecom product service url, example: http://127.0.0.1:8020>
```

6. create postgresql databases with name same as in .env file
7. test server first with `go test ./...` to make sure server works fine
8. run server with `go run ./...`
