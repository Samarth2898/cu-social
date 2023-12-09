package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	username  = "postgres"
	dbname    = "cusocial"
	pgaddress = "localhost"
)

type Post struct {
	PostID      int    `json:"post_id" db:"post_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	VideoURL    string `json:"video_url" db:"video_url"`
	UserID      int    `json:"user_id" db:"user_id"`
	Status      string `json:"status" db:"status"`
	CreatedAt   string `json:"created_at" db:"created_at"`
}

func (server *Server) getUserPosts(ctx *gin.Context) {

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", username, pgaddress, dbname))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var posts []Post

	// Query posts of the specified user from the database
	rows, err := db.Query("SELECT * FROM Posts WHERE user_id = $1", userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	fmt.Printf("rows %+v", rows)
	// Iterate through the rows and append posts to the slice
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Description, &post.VideoURL, &post.UserID, &post.Status, &post.CreatedAt); err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post row"})
			return
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, posts)
}

func (server *Server) getPosts(ctx *gin.Context) {

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", username, pgaddress, dbname))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	userIDStr := ctx.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var posts []Post

	// Query posts from users that the specified user is following
	query := `
		SELECT p.*
		FROM Posts p
		INNER JOIN follows f ON p.user_id = f.following_user_id
		WHERE f.followed_user_id = $1
		ORDER BY p.created_at DESC
		`

	rows, err := db.Query(query, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}
	defer rows.Close()

	// Iterate through the rows and append posts to the slice
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Description, &post.VideoURL, &post.UserID, &post.Status, &post.CreatedAt); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post row"})
			return
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, posts)
}

func (server *Server) createPost(ctx *gin.Context) {

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", username, pgaddress, dbname))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	var post Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(post)

	// Validate UserID exists in the users table (you might have a similar getUserID function)
	userIDExists := true // Replace this with logic to check if the UserID exists

	if !userIDExists {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}

	// Insert the post into the database
	insertPostQuery := `
			INSERT INTO posts (title, description, video_url, user_id, status)
			VALUES ($1, $2, $3, $4, $5)
		`

	fmt.Println(insertPostQuery, post.Title, post.Description, post.VideoURL, post.UserID, post.Status)

	// fmt.Println("Query:", queryString)

	_, err = db.Exec(insertPostQuery, post.Title, post.Description, post.VideoURL, post.UserID, post.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}
