package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
