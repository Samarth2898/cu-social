package main

import (
	"context"
	helper "cu-social/helper"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

// Book ...
type Book struct {
	Title  string
	Author string
}

type FeedObject struct {
	Description string `json:"description" db:"description"`
	VideoURL    string `json:"video_url" db:"video_url"`
	Title       string `json:"title" db:"title"`
	UserID      int    `json:"user_id" db:"user_id"`
	PostedBy    string `json:"username" db:"username"`
}

const (
	backendServer = "http://localhost:8080"
)

func main() {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "sa.json")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/styles", "./styles")
	r.Static("/js", "./js")
	r.Static("/images", "./images")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	r.POST("/submit-signup-form", func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		reqBody := struct {
			Username string
			Email    string
			Password string
		}{
			Username: username,
			Email:    email,
			Password: password,
		}
		marshalledBody, _ := json.Marshal(reqBody)
		res, err := helper.PostReq("http://0.0.0.0:8080/users", marshalledBody)
		if err != nil {
			fmt.Println("error sending POST request: ", err.Error())
		}
		fmt.Println(res.Status)
		c.HTML(http.StatusOK, "create_user_profile.html", gin.H{})
	})

	r.GET("/feed", func(c *gin.Context) {
		userID := 4 // samy get userID from JWT token
		FeedObjectInstance := getFeedByUser(userID)
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"FeedObjects": FeedObjectInstance,
		})
	})

	r.GET("/add_post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_post.html", gin.H{})
	})

	r.GET("/profile", func(c *gin.Context) {
		userID := 4 // samy get userID from JWT token
		FeedObjectInstance := getProfileFeed(userID)
		c.HTML(http.StatusOK, "profile.html", gin.H{
			"PostObjects": FeedObjectInstance,
		})
	})

	r.POST("/uploadVideo", func(c *gin.Context) {
		uploadVideoFunc(c)
	})

	log.Fatal(r.Run(":3000"))
}

func getFeedByUser(userID int) *[]FeedObject {

	url := fmt.Sprintf("%s/users/%d/following/posts", backendServer, userID)

	// Send GET request to fetch user posts
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching posts:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch posts. Status code: %d", response.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Unmarshal JSON response into Post slice
	var posts []FeedObject
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	return &posts

}

func getProfileFeed(userID int) *[]FeedObject {

	url := fmt.Sprintf("%s/users/%d/posts", backendServer, userID)

	// Send GET request to fetch user posts
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching posts:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch posts. Status code: %d", response.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// Unmarshal JSON response into Post slice
	var posts []FeedObject
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}

	return &posts

}

func uploadVideoFunc(c *gin.Context) {

	bucketName := "cusocial"

	ctx := context.Background()

	file, handler, err := c.Request.FormFile("videoFile")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file")
		return
	}
	defer file.Close()

	client, err := storage.NewClient(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to create GCS client")
		return
	}
	defer client.Close()

	objectName := handler.Filename
	obj := client.Bucket(bucketName).Object(objectName)
	wObj := obj.NewWriter(ctx)
	defer wObj.Close()

	if _, err := io.Copy(wObj, file); err != nil {
		c.String(http.StatusInternalServerError, "Error uploading file to GCS")
		return
	}

	if err := wObj.Close(); err != nil {
		c.String(http.StatusInternalServerError, "Error closing GCS writer")
		return
	}

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	c.String(http.StatusOK, "File uploaded successfully! URL: %s", url)
	c.Redirect(http.StatusMovedPermanently, "/feed")
}
