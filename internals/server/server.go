package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sub-rat/machnet_api_assingment/internals/features/post"
	"github.com/sub-rat/machnet_api_assingment/internals/features/timeline"
	"github.com/sub-rat/machnet_api_assingment/internals/features/user"
	"github.com/sub-rat/machnet_api_assingment/pkg/db/postgres"
	"gorm.io/gorm"
)

type server struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func GetServer() *server {
	return &server{
		Engine: gin.Default(),
		DB:     postgres.ConnectDatabase(),
	}
}

func (server *server) Run() {
	server.initRoutes()
	log.Fatal(server.Engine.Run())
}

func (server *server) initRoutes() {
	// routes or Endpoints
	eng := server.Engine
	user.RegisterRoutes(eng, user.NewService(user.NewRepository(*server.DB)))
	post.RegisterRoutes(eng, post.NewService(post.NewRepository(*server.DB), timeline.NewRepository(*server.DB)))
	timeline.RegisterRoutes(eng, timeline.NewService(timeline.NewRepository(*server.DB), user.NewRepository(*server.DB)))

}
