package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Post represents the structure of a post in the database
type Post struct {
	PostID          int    `db:"post_id"`
	Title			string `db:"title"`
	Description		string `db:"description"`
	VideoUrl		string `db:"video_url"`
	UserID          int    `db:"user_id"`
	Status          string `db:"status"`
	CreatedAt		string `db:"created_at"`
}

CREATE TABLE "posts" (
	"post_id" SERIAL PRIMARY KEY,
	"title" varchar,
	"description" varchar,
	"video_url" text,
	"user_id" integer,
	"status" varchar,
	"created_at" timestamp
  );

func main() {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "postgres://username:password@localhost/database_name?sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Initialize the Gin router
	r := gin.Default()

	// Endpoint to get all posts of a user
	r.GET("/users/:userID/posts", func(c *gin.Context) {
		userIDStr := c.Param("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var posts []Post

		// Query posts of the specified user from the database
		rows, err := db.Query("SELECT * FROM Posts WHERE UserID = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		defer rows.Close()

		// Iterate through the rows and append posts to the slice
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.ImageOrVideoURL); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post row"})
				return
			}
			posts = append(posts, post)
		}

		c.JSON(http.StatusOK, posts)
	})

	// Run the server
	r.Run(":8080")
}
