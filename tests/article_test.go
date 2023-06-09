package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteArticle(t *testing.T) {
	slug, createdUserToken := CreateArticle()
	req, err := http.NewRequest("DELETE",host+"/api/articles/"+slug, nil)

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

	// Initialize TagList as an empty slice in the struct
	payload := ArticleDeleteResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err) 
		fmt.Println(" Variable: ", payload)
	}

	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Response Body:", string(body))

	assert.Equal(t, 200, resp.StatusCode, "Expected status code to be 200")
	assert.Equal(t, "200 OK", resp.Status, "Expected status to be 200 OK")
	assert.Equal(t, "Delete success", payload.Article, "Expected status to be Delete success")

	reqGetArticle, err := http.NewRequest("GET",host+"/api/articles/"+slug, nil)

	// Set headers
	reqGetArticle.Header.Set("Content-Type", "application/json")
	reqGetArticle.Header.Set("Authorization", "Token "+token)

	getClient := http.DefaultClient
	respGetArticle, err := getClient.Do(reqGetArticle)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer respGetArticle.Body.Close()

	bodyGetArticle, errGetArticle := ioutil.ReadAll(respGetArticle.Body)
	if errGetArticle != nil {
		fmt.Println("Error reading response body:", errGetArticle.Error())
		return
	}

	// Initialize TagList as an empty slice in the struct
	payloadGetArticle := ArticleResponse{}
    err = json.NewDecoder(bytes.NewReader(bodyGetArticle)).Decode(&payloadGetArticle)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err) 
		fmt.Println(" Variable: ", payload)
	}

	assert.Equal(t, "", payloadGetArticle.Article.Slug, "Expected slug to be empty")
}

func TestTryDeleteArticleFromAnotherUser(t *testing.T) {
	slug, userTokenOwnerArticle := CreateArticle()
	fmt.Println("createdUserTokenWithArticle:", userTokenOwnerArticle)
	req, err := http.NewRequest("DELETE",host+"/api/articles/"+slug, nil)

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token Ababuble")

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

	reqGetArticle, err := http.NewRequest("GET",host+"/api/articles/"+slug, nil)

	// Set headers
	reqGetArticle.Header.Set("Content-Type", "application/json")
	reqGetArticle.Header.Set("Authorization", "Token "+token)

	getClient := http.DefaultClient
	respGetArticle, err := getClient.Do(reqGetArticle)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer respGetArticle.Body.Close()

	body, err := ioutil.ReadAll(respGetArticle.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err.Error())
		return
	}

	// Initialize TagList as an empty slice in the struct
	payload := ArticleResponse{}
    err = json.NewDecoder(bytes.NewReader(body)).Decode(&payload)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err) 
		fmt.Println(" Variable: ", payload)
	}

	assert.Equal(t, slug, payload.Article.Slug, "Expected slug to exists")
}