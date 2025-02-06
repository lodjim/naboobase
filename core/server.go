package core

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
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
			server.Router.POST(endpoint.Path, endpoint.Handler)
		case "GET":
			server.Router.GET(endpoint.Path, endpoint.Handler)
		case "PATCH":
			server.Router.PATCH(endpoint.Path, endpoint.Handler)
		case "PUT":
			server.Router.PUT(endpoint.Path, endpoint.Handler)
		case "DELETE":
			server.Router.DELETE(endpoint.Path, endpoint.Handler)
		case "OPTIONS":
			server.Router.OPTIONS(endpoint.Path, endpoint.Handler)
		}
	}
}

func (server *Server) AttachAuthenticationLayer(db MongoDBconnector) {
	auth := Authenticator{
		Type: PasswordTypeAuthenticator,
	}
	server.AttachEndpoints(auth.Init(db))
}

func (server *Server) AutoServe(db MongoDBconnector) {
	var newEndpoints []Endpoint
	for k, v := range AutoEndpointFuncRegistry {
		newEndpoint := Endpoint{
			Method:  "POST",
			Handler: v(db),
			Path:    fmt.Sprintf("/%s", k),
		}
		newEndpoints = append(newEndpoints, newEndpoint)
	}
	server.AttachEndpoints(newEndpoints)
}

func (server *Server) RunServer() {
	link := fmt.Sprintf("%s:%v", server.IP, server.Port)
	err := server.Router.Run(link)
	if err != nil {
		fmt.Println(err.Error())
	}
}
