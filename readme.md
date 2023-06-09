# RealWorld Example App test with Go

This codebase was created to demonstrate a fully API built with **Golang/Gin** with tests using GO.
The backend test was done using just Go, testing an application built in GO that is a medium.com backend clone. You can find the repository here https://github.com/gothinkster/golang-gin-realworld-example-app

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

## Troubleshooting

- Undefined validation function 'exists' on field 'Username'
  If you get this error you should try downgrade to gin-gonic/gin@v1.4.0 using:

```
go get github.com/gin-gonic/gin@v1.4.0
```

Then build and execute the project like:

```
go build
./golang-gin-realworld-example-app
```

### How to run the tests in GHA

With the tests automated in the GHA pipeline, we can execute our tests every time a pull request is opened to our project to assure everything still works like expected and block the Pull Request from being merged if any tests failed

We can also execute the tests manually by the workflow dispatch in the following link by clicking the [following link](https://github.com/CaiqueCoelho/go-test/actions/workflows/api-tests.yml) on Run Workflow in the gray box and on the Run Workflow green button

https://github.com/CaiqueCoelho/go-test/actions/workflows/api-tests.yml

![Captura de Tela 2023-06-09 às 17.10.39.png](https://s3-us-west-2.amazonaws.com/secure.notion-static.com/7f5d888c-ea36-4dca-93a7-bd035e80a084/Captura_de_Tela_2023-06-09_as_17.10.39.png)

You can the result of the tests https://github.com/CaiqueCoelho/go-test/actions/runs/5225995752/jobs/9436167002

And the JSON report in the following link https://github.com/CaiqueCoelho/go-test/suites/13502035308/artifacts/741815992

### How to run the project locally and run the tests

Install Dependencies and run the project

1.  **From the project root, run in the terminal:**

```jsx
go build
go mod tidy
./golang-gin-realworld-example-app
```

1.  **From the /tests path run:**

`go test -v`

### About the tests

To test our backend from a user perspective we have lunch the application running the application locally in localhost:8080.

All tests were included in a directory called tests and in a file called users_test.

About the architecture and organization of the tests:

We have added all the tests in a directory called tests, s**o to run the tests you should be in the path** `/tests`.

We have three files:

1. main_tests.go → Main entry to call all the tests in the other two files, this file was used to centralize and abstract all the auxiliary functions, variables, types, and structs to be used in our tests, for example, methods to Create Article and User to be used in our tests for example to tests an update in a user we should have a created user first, so to avoid duplicated code I’ve created those auxiliary functions
2. users_tests.go → With all tests in the user scope
3. article_test.go → With all tests in the article scope

The tests files were added to a package called main_test, in this test file we also added the gofakeit and testify packages to help us generate fake random data to be used, like emails, and also to help us do assertions validation in the responses from our requests tests.

I’ve covered 4 flows in the backend activity:

- POST api/users
  - A successful test
    - To test this I’ve sent a post request to [http://localhost:8080](http://localhost:8080/)/api/users with an email random generated from the gofakeit package
    - I’ve made assertions to check if:
      - The status code and status be 201 Created
      - The username returned in the response is the same one sent in the request body
      - The Bio returned in the response is empty because wasn’t sent in the request body
      - The image returned in the response is null because wasn’t sent in the request body
      - The username email in the response is the same one sent in the request body
      - The token returned in the response is not empty
  - A test with an empty email field
    - I’ve made assertions to check if:
      - The status code and status be 422 Unprocessable Entity
      - The return error in the body is `{key: email}`
  - A test validation email unique constraint
    - I’ve made assertions to check if:
      - The status code and status be 422 Unprocessable Entity
      - The return error in the body is `UNIQUE constraint failed: user_models.email`
- GET api/user
  - Successful test
    - To test this I’ve sent a get request to http://localhost:8080/api/user with the token returned from the first test when the user is created
    - I’ve made assertions to check if:
      - The status code and status be 200 OK
      - The username returned in the response is the same one sent in the auxiliary method CreateUser()
      - The Bio returned in the response is empty because wasn’t sent in the request body
      - The image returned in the response is null because wasn’t sent in the request body
      - The username email in the response is the same one sent in the auxiliary method CreateUser()
      - The token returned is the same one sent in the auxiliary method CreateUser()
  - Getting a user without token
    - I’ve made assertions to check if:
      - The status code and status be 401 Unauthorized
- PUT /api/user
  - Edit an user with success:
    - I’ve made assertions to check if:
      - The status code and status be 200 OK
      - The username returned in the response is the same one sent in the auxiliary method CreateUser()
      - The email returned in the response is the same one sent in the auxiliary method CreateUser()
      - The token returned in the response is the same one sent in the auxiliary method CreateUser()
  - Try to edit letting the e-mail empty
    - I’ve made assertions to check if:
      - The status code and status be 422 Unprocessable Entity
      - The return error in the body is `{key: email}`
  - Try to edit using an invalid e-mail
    - I’ve made assertions to check if:
      - The status code and status be 422 Unprocessable Entity
      - The return error in the body is `{key: email}`
  - Try to edit using an already-registered email
    - I’ve made assertions to check if:
      - The status code and status be 422 Unprocessable Entity
      - The return error in the body is `UNIQUE constraint failed: user_models.email`
- DELETE /api/articles/:slug
  - Delete an article with slug with success:
    - I’ve made assertions to check if:
      - The status code and status be 200 OK
      - The message returned is Delete success
      - The slug article doesn’t exists
  - Try to delete an article with an invalid token:
    - I’ve made assertions to check if:
      - The status code and status be 401 Unauthorized
      - The username returned in the response is the same one sent in the auxiliary method CreateUser()
      - The slug article still exists

### About bugs and problems found

1. Try to delete an article with a non-exist slug or an article with a slug from another user:
   1. Description: When trying to delete an article with a non-exist slug or passing a slug of another user article the API always returns 200 OK with the message Delete success, if the slug is from a non-exist article nothing happens, if the slug is from another user different from the token passed in the authorization the article isn’t delete
   2. Steps to reproduce:
      1. Call the DELETE method in the URL http://localhost:8080/api/articles/non-exist-slug
      2. Make sure you don’t have a published article with the slug non-exist-slug
   3. Obtained result: 200 OK
   4. Expect result: 404 Not found or another status that tells us the delete wasn’t a success because the slug doesn’t exist
2. **Undefined validation function 'exists' on field 'Username'**
   1. The first difficulty I found ins this repository was about how to run the project since when I executed and was doing a POST request to /api/users I was getting the error **`Undefined validation function 'exists' on field 'Username'`** after some investigations I discovered that was something about the gin version so I downgrade to the `gin@v1.4.0` version and I was able to do the POST request, so I have created an [Issue](https://github.com/gothinkster/golang-gin-realworld-example-app/issues/32) in the repository about that problem and how to fix this and I also open a [Pull Request](https://github.com/gothinkster/golang-gin-realworld-example-app/pull/33) with a Troubleshooting session and a better description about how to run the project
   2. Issue https://github.com/gothinkster/golang-gin-realworld-example-app/issues/32
   3. Pull Request: https://github.com/gothinkster/golang-gin-realworld-example-app/pull/33
