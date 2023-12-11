package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "csbackend/db/sqlc"
	"csbackend/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type User struct {
	UserID         int    `db:"user_id" json:"user_id"`
	Username       string `db:"username" json:"username"`
	Password       string `db:"password" json:"password"`
	ProfilePicture string `db:"profile_picture" json:"profile_picture"`
	Biography      string `db:"biography" json:"biography"`
	Email          string `db:"email" json:"email"`
	CreatedAt      string `db:"created_at" json:"created_at"`
}

type createUserRequest struct {
	Username       string `json:"username" binding:"required,alphanum"`
	Password       string `json:"password" binding:"required,min=6"`
	Email          string `json:"email" binding:"required,email"`
	ProfilePicture string `json:"profile_picture"`
	Biography      string `json:"biography"`
}

type userResponse struct {
	UserId    int32     `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserId:    user.UserID,
		Username:  user.Username.String,
		Email:     user.Email.String,
		CreatedAt: user.CreatedAt.Time,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	arg := db.CreateUserParams{
		Username: sql.NullString{String: req.Username, Valid: true},
		Password: sql.NullString{String: hashedPassword, Valid: true},
		Email:    sql.NullString{String: req.Email, Valid: true},
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, sql.NullString{String: req.Username, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password.String)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Username.String,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)

}

type UserData struct {
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type searchUserResponse struct {
	UserData []UserData `json:"user_data"`
}

// func (server *Server) searchUsers(ctx *gin.Context) {
// 	userID, _ := strconv.Atoi(ctx.Query("user_id"))
// 	searchQuery := ctx.Query("query") + "%"
// 	queryReq := db.SearchUsersParams{
// 		UserID:   int32(userID),
// 		Username: sql.NullString{String: searchQuery, Valid: true},
// 	}
// 	rows, err := server.store.SearchUsers(ctx, queryReq)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errResponse(err))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errResponse(err))
// 		return
// 	}
// 	var userInfo []UserData

// 	for _, row := range rows {
// 		var currentUser UserData
// 		currentUser.Username = row.Username.String
// 		currentUser.ProfilePicture = row.ProfilePicture.String
// 		userInfo = append(userInfo, currentUser)
// 	}

// 	rsp := searchUserResponse{
// 		UserData: userInfo,
// 	}
// 	ctx.JSON(http.StatusOK, rsp)
// }

func (server *Server) searchUsers(ctx *gin.Context) {
	db, err := sql.Open(server.config.DBDriver, server.config.DBSource)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	var users []User

	query := "SELECT * FROM users"
	rows, err := db.Query(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Username, &user.Password, &user.ProfilePicture, &user.Biography, &user.Email, &user.CreatedAt)
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan users"})
			return
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over users"})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

type updateUserDataRequest struct {
	ProfilePicture string `json:"profile_picture"`
	Biography      string `json:"biography"`
	UserID         int32  `json:"user_id"`
}

type updateUserDataResponse struct {
	Success bool `json:"success"`
}

func (server *Server) updateUserInfo(ctx *gin.Context) {
	var req updateUserDataRequest
	var dbReq db.UpdateUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	dbReq = db.UpdateUserParams{
		UserID:         req.UserID,
		ProfilePicture: sql.NullString{String: req.ProfilePicture, Valid: true},
		Biography:      sql.NullString{String: req.Biography, Valid: true},
	}
	fmt.Println("inside update: ", dbReq)
	res, err := server.store.UpdateUser(ctx, dbReq)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	var status bool
	if res == 1 {
		status = true
	} else {
		status = false
	}
	rsp := updateUserDataResponse{
		Success: status,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) userInfo(ctx *gin.Context) {
	db, err := sql.Open(server.config.DBDriver, server.config.DBSource)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	id := ctx.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Query to fetch user info
	query := "SELECT user_id, username, profile_picture, biography, email, created_at FROM users WHERE user_id = $1"
	var user User
	err = db.QueryRow(query, userID).Scan(&user.UserID, &user.Username, &user.ProfilePicture, &user.Biography, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
