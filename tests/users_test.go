package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Log("Host:", host)
	t.Log("Email:", email)
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(createUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	resp, err := http.Post(host+"/api/users", "application/json", buffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := UserResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	token = string(payload.User.Token)

	assert.Equal(t, 201, resp.StatusCode, "Expected status code to be 201")
	assert.Equal(t, "201 Created", resp.Status, "Expected status to be 201 Created")
	assert.Equal(t, createUserRequestBody.User.Username, payload.User.Username, "Expected username to be " + createUserRequestBody.User.Username)
	assert.Equal(t, "", payload.User.Bio, "Expected Bio to be an empty string")
	assert.Nil(t, payload.User.Image, "Expected image to be null")
	assert.Equal(t, createUserRequestBody.User.Email, payload.User.Email, "Expected email to be" + createUserRequestBody.User.Email)
	assert.NotEmpty(t, payload.User.Token, "Expected token to be not empty")
}

func TestCreateUserWithAlreadyRegisteredEmail(t *testing.T) {
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(createUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	resp, err := http.Post(host+"/api/users", "application/json", buffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := ErrorResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 422, resp.StatusCode, "Expected status code to be 422")
	assert.Equal(t, "422 Unprocessable Entity", resp.Status, "Expected status to be 422 Unprocessable Entity")
	assert.Equal(t, "UNIQUE constraint failed: user_models.email", payload.Errors.Database, "Expected error to be UNIQUE constraint failed: user_models.email")
}

func TestCreateUserWithEmptyEmail(t *testing.T) {
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	requestBodyWithEmptyEmail := createUserRequestBody
	requestBodyWithEmptyEmail.User.Email = ""

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(requestBodyWithEmptyEmail); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	resp, err := http.Post(host+"/api/users", "application/json", buffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := ErrorResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 422, resp.StatusCode, "Expected status code to be 422")
	assert.Equal(t, "422 Unprocessable Entity", resp.Status, "Expected status to be 422 Unprocessable Entity")
	assert.Equal(t, "{key: email}", payload.Errors.Email, "Expected error to be key: email")
}

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET",host+"/api/user", nil)

	createdUser := CreateUser()
	createdUserToken := string(createdUser.User.Token)
	createdUserEmail := createdUser.User.Email
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+createdUserToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := UserResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Equal(t, "200 OK", resp.Status, "Expected status to be 200 OK")
	assert.Equal(t, createUserRequestBody.User.Username, payload.User.Username, "Expected username to be " + createUserRequestBody.User.Username)
	assert.Equal(t, "", payload.User.Bio, "Expected Bio to be an empty string")
	assert.Nil(t, payload.User.Image, "Expected image to be null")
	assert.Equal(t, createdUserEmail, payload.User.Email, "Expected email to be" + createUserRequestBody.User.Email)
	assert.Equal(t, createdUserToken, payload.User.Token, "Expected token to be the same one send in header authorization")
}


func TestGetUserWithouTokenAuth(t *testing.T) {
	req, err := http.NewRequest("GET",host+"/api/user", nil)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)

	assert.Equal(t, 401, resp.StatusCode, "Expected status code to be 401")
	assert.Equal(t, "401 Unauthorized", resp.Status, "Expected status to be 401 Unauthorized")
}

func TestUpdateUserEmail(t *testing.T) {
	updateUserRequestBody := UserUpdate{
		User: UserUpdatePayload{
			Email:    updateEmail,
		},
	}
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(updateUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT",host+"/api/user", buffer)

	createdUser := CreateUser()
	createdUserToken := string(createdUser.User.Token)
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+createdUserToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := UserResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Equal(t, "200 OK", resp.Status, "Expected status to be 200 OK")
	assert.Equal(t, createUserRequestBody.User.Username, payload.User.Username, "Expected username to be " + createUserRequestBody.User.Username)
	assert.Equal(t, updateEmail, payload.User.Email, "Expected email to be" + createUserRequestBody.User.Email)
	assert.Equal(t, createdUserToken, payload.User.Token, "Expected token to be the same one send in header authorization")
}

func TestUpdateUserWithInvalidEmail(t *testing.T) {
	updateUserRequestBody := UserUpdate{
		User: UserUpdatePayload{
			Email:    "caique",
		},
	}
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(updateUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT",host+"/api/user", buffer)
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := ErrorResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 422, resp.StatusCode, "Expected status code to be 422")
	assert.Equal(t, "422 Unprocessable Entity", resp.Status, "Expected status to be 422 Unprocessable Entity")
	assert.Equal(t, "{key: email}", payload.Errors.Email, "Expected username to be {key: email}")
}

func TestUpdateUserWithAlreadyRegisteredEmail(t *testing.T) {
	// Create a bytes.Buffer
	createUserRequestbuffer := bytes.NewBuffer(nil)
	alreadyRegistedUserRequestBody := createUserRequestBody
	alreadyRegistedEmail := gofakeit.Email()
	alreadyRegistedUserRequestBody.User.Email = alreadyRegistedEmail

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(createUserRequestbuffer).Encode(alreadyRegistedUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	respCreatedUser, err := http.Post(host+"/api/users", "application/json", createUserRequestbuffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer respCreatedUser.Body.Close()

	updateUserRequestBody := UserUpdate{
		User: UserUpdatePayload{
			Email:    alreadyRegistedEmail,
		},
	}
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(updateUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT",host+"/api/user", buffer)
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := ErrorResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 422, resp.StatusCode, "Expected status code to be 422")
	assert.Equal(t, "422 Unprocessable Entity", resp.Status, "Expected status to be 422 Unprocessable Entity")
	assert.Equal(t, "UNIQUE constraint failed: user_models.email", payload.Errors.Database, "Expected username to be UNIQUE constraint failed: user_models.email")
}


func TestUpdateUserWithEmptyEmail(t *testing.T) {
	// Create a bytes.Buffer
	createUserRequestbuffer := bytes.NewBuffer(nil)
	alreadyRegistedUserRequestBody := createUserRequestBody
	newEmail := gofakeit.Email()
	alreadyRegistedUserRequestBody.User.Email = newEmail

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(createUserRequestbuffer).Encode(alreadyRegistedUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	respCreatedUser, err := http.Post(host+"/api/users", "application/json", createUserRequestbuffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer respCreatedUser.Body.Close()

	updateUserRequestBody := UserUpdate{
		User: UserUpdatePayload{
			Email:    "",
		},
	}
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(updateUserRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT",host+"/api/user", buffer)
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	payload := ErrorResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 422, resp.StatusCode, "Expected status code to be 422")
	assert.Equal(t, "422 Unprocessable Entity", resp.Status, "Expected status to be 422 Unprocessable Entity")
	assert.Equal(t, "{key: email}", payload.Errors.Email, "Expected username to be UNIQUE constraint failed: user_models.email")
}