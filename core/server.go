package core

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	IP     string
	Port   int
	Router *gin.Engine
}

func (server *Server) Init(ip string, port int) {
	server.IP = ip
	server.Port = port
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

func (server *Server) AttachEndpoints(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		switch endpoint.Method {
		case "POST":
			server.Router.POST(endpoint.Path, *endpoint.Handler)
		case "GET":
			server.Router.GET(endpoint.Path, *endpoint.Handler)
		case "PATCH":
			server.Router.PATCH(endpoint.Path, *endpoint.Handler)
		case "PUT":
			server.Router.PUT(endpoint.Path, *endpoint.Handler)
		case "DELETE":
			server.Router.DELETE(endpoint.Path, *endpoint.Handler)
		case "OPTIONS":
			server.Router.OPTIONS(endpoint.Path, *endpoint.Handler)
		}
	}
}

func (server *Server) RunServer() {
	link := fmt.Sprintf("http://%s:%v", server.IP, server.Port)
	err := server.Router.Run(link)
	if err != nil {
		fmt.Println(err.Error())
	}
}
