package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Follow struct {
	FollowingUserID int `json:"following_user_id"`
	FollowedUserID  int `json:"followed_user_id"`
}

func (server *Server) followUser(ctx *gin.Context) {

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", username, pgaddress, dbname))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	followerID := ctx.Param("followerID")
	followingID := ctx.Param("followingID")

	var follow Follow
	err = ctx.BindJSON(&follow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add logic to check if the followerID and followingID exist in the users table -- Sammy
	// Assuming followerID and followingID are valid user IDs

	// Insert follow relationship into the follows table
	insertFollowQuery := "INSERT INTO follows (following_user_id, followed_user_id) VALUES ($1, $2)"
	_, err = db.Exec(insertFollowQuery, followerID, followingID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to follow user may be redundant request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

func (server *Server) checkFollowUser(ctx *gin.Context) {

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", username, pgaddress, dbname))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	followingID := ctx.Param("followingID")
	followedID := ctx.Param("followedID")

	query := "SELECT COUNT(*) FROM follows WHERE following_user_id = $1 AND followed_user_id = $2"
	var count int
	err = db.QueryRow(query, followingID, followedID).Scan(&count)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Database error"})
		return
	}

	if count > 0 {
		ctx.JSON(200, gin.H{"is_following": true})
	} else {
		ctx.JSON(200, gin.H{"is_following": false})
	}
}
