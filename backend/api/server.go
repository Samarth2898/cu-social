package api

import (
	"fmt"

	db "csbackend/db/sqlc"
	"csbackend/token"
	"csbackend/util"

	"github.com/gin-gonic/gin"
)

// serve HTTP requests for banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.GET("/users/list", server.searchUsers)
	router.PUT("/users/update", server.updateUserInfo)
	router.GET("/users/:userID/posts", server.getUserPosts)
	router.GET("/users/:userID/following/posts", server.getPosts)
	// router.POST("/users/:followerID/follow/:followingID", server.followUser)
	router.POST("/posts", server.createPost)

	fmt.Println(authRoutes)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
