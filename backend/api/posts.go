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
	username = "postgres"
	dbname   = "cusocial"
)

// Post represents the structure of a post in the database
type Post struct {
	PostID      int    `db:"post_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	VideoUrl    string `db:"video_url"`
	UserID      int    `db:"user_id"`
	Status      string `db:"status"`
	CreatedAt   string `db:"created_at"`
}

func (server *Server) getPosts(ctx *gin.Context) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@localhost/%s?sslmode=disable", username, dbname))
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
		if err := rows.Scan(&post.PostID, &post.Title, &post.Description, &post.VideoUrl, &post.UserID, &post.Status, &post.CreatedAt); err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan post row"})
			return
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, posts)
}

// func (server *Server) getPosts(ctx *gin.Context) {
// 	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@localhost/%s?sslmode=disable", username, dbname))
// 	if err != nil {
// 		fmt.Println("Error connecting to the database:", err)
// 		return
// 	}

// 	query := "SELECT * FROM posts"

// 	// Print the SQL query
// 	fmt.Println("SQL Query:", query)

// 	// Execute the query
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	// Iterate through the rows and print each row's columns
// 	for rows.Next() {
// 		columns, err := rows.Columns()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		values := make([]interface{}, len(columns))
// 		scanArgs := make([]interface{}, len(columns))

// 		for i := range values {
// 			scanArgs[i] = &values[i]
// 		}

// 		if err := rows.Scan(scanArgs...); err != nil {
// 			log.Fatal(err)
// 		}

// 		fmt.Println("Row Data:")
// 		for i, col := range values {
// 			fmt.Printf("%s: %v\n", columns[i], col)
// 		}
// 	}

// 	// Check for errors during row iteration
// 	if err := rows.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// }
