package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	db "csbackend/db/sqlc"
	"csbackend/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

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

func (server *Server) searchUsers(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))
	searchQuery := ctx.Query("query") + "%"
	queryReq := db.SearchUsersParams{
		UserID:   int32(userID),
		Username: sql.NullString{String: searchQuery, Valid: true},
	}
	rows, err := server.store.SearchUsers(ctx, queryReq)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	var userInfo []UserData

	for _, row := range rows {
		var currentUser UserData
		currentUser.Username = row.Username.String
		currentUser.ProfilePicture = row.ProfilePicture.String
		userInfo = append(userInfo, currentUser)
	}

	rsp := searchUserResponse{
		UserData: userInfo,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type updateUserDataResponse struct {
	Success bool `json:"success"`
}

func (server *Server) updateUserInfo(ctx *gin.Context) {
	var req db.UpdateUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	res, err := server.store.UpdateUser(ctx, req)
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
