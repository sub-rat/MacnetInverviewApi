package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sub-rat/social_network_api/internals/features/post"
	"github.com/sub-rat/social_network_api/internals/features/timeline"
	"github.com/sub-rat/social_network_api/internals/features/user"
	"github.com/sub-rat/social_network_api/pkg/db/postgres"
	"gorm.io/gorm"
)

type server struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

// GetServer Returns the server instance with postgres connection and Gin initialization
func GetServer() *server {
	return &server{
		Engine: gin.Default(),
		DB:     postgres.ConnectDatabase(),
	}
}

// Run initialize the routes and start the server engine
func (server *server) Run() {
	server.initRoutes()
	log.Fatal(server.Engine.Run())
}

// initRoutes initializes the Enginer and Routes
func (server *server) initRoutes() {
	eng := server.Engine
	user.RegisterRoutes(eng, user.NewService(user.NewRepository(*server.DB)))
	post.RegisterRoutes(eng, post.NewService(post.NewRepository(*server.DB), timeline.NewRepository(*server.DB)))
	timeline.RegisterRoutes(eng, timeline.NewService(timeline.NewRepository(*server.DB), user.NewRepository(*server.DB)))

}
