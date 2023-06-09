package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

var (
	email string
	updateEmail string
	token string
	host string
	createUserRequestBody UserCreateRequest
	createArticleRequestBody ArticleCreateRequest
	title string
)

type UserResponse struct {
	User UserResponsePayload `json:"user"`
}

type UserUpdate struct {
	User UserUpdatePayload `json:"user"`
}

type UserUpdatePayload struct {
	Email string `json:"email"`
}

type ErrorResponse struct {
	Errors ErrorEmailResponse `json:"errors"`
}

type ErrorEmailResponse struct {
	Email string `json:"Email"`
	Database string `json:"database"`
}

type UserResponsePayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    interface{} `json:"image"`
	Token    string `json:"token"`
}

type UserCreateRequest struct {
	User UserRequest `json:"user"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type ArticleCreateRequest struct {
	Article ArticleRequest `json:"article"`
}

type ArticleRequest struct {
	Title    string `json:"title"`
	Description string `json:"description"`
	Body string `json:"body"`
	TagList []string `json:"tagList"`
}

type ArticleResponse struct {
	Article ArticleBody `json:"article"`
}

type ArticleBody struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	Body          string `json:"body"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	Author        AuthorBody `json:"author"`
	TagList        []string `json:"tagList"`
	Favorited      bool     `json:"favorited"`
	FavoritesCount int      `json:"favoritesCount"`
}

type AuthorBody struct {
	Username  string  `json:"username"`
	Bio       string  `json:"bio"`
	Image     **string `json:"image"`
	Following bool    `json:"following"`
}

type ArticleDeleteResponse struct {
	Article string    `json:"article"`
}

func CreateUser() UserResponse {
	// Create the request body
	createUser := UserCreateRequest{
		User: UserRequest{
			Email:    gofakeit.Email(),
			Password: "test123456",
			Username: "caiquecoelho15",
		},
	}
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(createUser); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	resp, err := http.Post(host+"/api/users", "application/json", buffer)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return UserResponse{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return UserResponse{}
	}

	payload := UserResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	return payload
}

// setup is called before running the tests.
func setup() {
	// Set the random seed for deterministic results (optional)
	gofakeit.Seed(0)
	email = gofakeit.Email()
	updateEmail = gofakeit.Email()
	host = "http://localhost:8080"
	// Create the request body
	createUserRequestBody = UserCreateRequest{
		User: UserRequest{
			Email:    email,
			Password: "test123456",
			Username: "caiquecoelho15",
		},
	}

	fmt.Println("Setup complete")
}

// TestMain is the entry point for running tests.
func TestMain(m *testing.M) {
	// Test setup
	setup()

	// Run the tests
	exitCode := m.Run()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

func CreateArticle() (string, string) {
	createdUser := CreateUser()
	createdUserToken := string(createdUser.User.Token)
	// Create a bytes.Buffer
	buffer := bytes.NewBuffer(nil)

	title = gofakeit.Sentence(7)

	createArticleRequestBody = ArticleCreateRequest{
		Article: ArticleRequest{
			Title:    title,
			Description: gofakeit.Sentence(10),
			Body: gofakeit.Sentence(30),
			TagList: []string{gofakeit.Word(), gofakeit.Word()},
		},
	}

	// Encode MyRequest directly into the buffer
	if err := json.NewEncoder(buffer).Encode(createArticleRequestBody); err != nil {
		// Handle the error
		log.Fatal(err)
	}

	// Send the POST request
	req, err := http.NewRequest("POST",host+"/api/articles", buffer)
	
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
		return "", ""
	}

	// Initialize TagList as an empty slice in the struct
	payload := ArticleResponse{
		Article: ArticleBody {
			TagList: []string{},
		},
	}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err) 
		fmt.Println(" Variable: ", payload)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	return payload.Article.Slug, createdUserToken
}