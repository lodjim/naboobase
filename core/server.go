package core

import (
	"fmt"
	"time"

	"github.com/aws/smithy-go/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	IP     string
	Port   int
	Router *gin.Engine
}

func (server *Server) Init(Ip string, port int) {

	server.Router = gin.Default()
	server.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                // Specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}, // Specify allowed methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Specify allowed headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // Allow credentials (cookies)
		MaxAge:           12 * time.Hour, // Maximum age for preflight request caching
	}))

	err := server.Router.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		fmt.Println(err.Error())
	}

}

func (server *Server) AttachMiddleware(middleware ...gin.HandlerFunc) {
	server.Router.Use(middleware...)
}
