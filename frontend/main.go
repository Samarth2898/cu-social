package main

import (
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
	Description string
	VideoURL    string
	Title       string
	PostedBy    string
}

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("templates/gin-test.html", "templates/feed.html")

	r.Static("/styles", "./styles")
	r.Static("/js", "./js")
	r.Static("/images", "./images")

	books := make([]Book, 0)
	books = append(books, Book{
		Title:  "Title 1",
		Author: "Author 1",
	})
	books = append(books, Book{
		Title:  "Title 2",
		Author: "Author 2",
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "gin-test.html", gin.H{
			"books": books,
		})
	})

	FeedObjectInstance := make([]FeedObject, 0)
	FeedObjectInstance = append(FeedObjectInstance, FeedObject{
		Description: "By Blender Foundation",
		VideoURL:    "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
		Title:       "Big Buck Bunny",
		PostedBy:    "Test User",
	})

	FeedObjectInstance = append(FeedObjectInstance, FeedObject{
		Description: "The first Blender Open Movie from 2006",
		VideoURL:    "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
		Title:       "Elephant Dream",
		PostedBy:    "Test User",
	})

	r.GET("/feed", func(c *gin.Context) {
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"FeedObjects": FeedObjectInstance,
		})
	})

	log.Fatal(r.Run())
}
