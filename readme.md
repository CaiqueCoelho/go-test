# RealWorld Example App test with Go

This codebase was created to demonstrate a fully API built with **Golang/Gin** with tests using GO

# Directory structure

```
.
├── gorm.db
├── hello.go
├── common
│   ├── utils.go        //small tools function
│   └── database.go     //DB connect manager
├── users
|   ├── models.go       //data models define & DB operation
|   ├── serializers.go  //response computing & format
|   ├── routers.go      //business logic & router binding
|   ├── middlewares.go  //put the before & after logic of handle request
|   └── validators.go   //form/json checker
├── tests
|   ├── main_tests.go       //main entry of all tests with auxiliary functions, variables, types and structs
|   ├── users_tests.go      //tests of the user scope
|   ├── article_tests.go    //tests of the article scope
...
```

# Getting started

## Install Golang

Make sure you have Go 1.13 or higher installed.

https://golang.org/doc/install

If you get the error `Undefined validation function 'exists' on field 'Username'`downgrade your gin-gonic/gin to 1.4.0
by doing `go get github.com/gin-gonic/gin@v1.4.0`

## Install Dependencies and running the project

From the project root, run:

```
go build
go mod tidy
./golang-gin-realworld-example-app
```

## Api Testing

From the /tests path run:

```
go test -v
```
