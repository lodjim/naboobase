package main

import (
	"net/http"

	"github.com/lodjim/naboobase/controllers"
	"github.com/lodjim/naboobase/core"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Good")
	}
}

var dbConnector = core.MongoDBconnector{}

func main() {
	_ = dbConnector.Connect("naboobase")
	myApi := core.Server{}
	myApi.Init("localhost", 1555)
	myApi.AttachEndpoints([]core.Endpoint{
		{
			Method:  "POST",
			Path:    "/user",
			Handler: controllers.CreateUser(dbConnector),
		},
		{
			Method:  "GET",
			Path:    "/health",
			Handler: HealthCheck(),
		},
	})
	myApi.AttachAuthenticationLayer(dbConnector)
	myApi.AutoServe(dbConnector)
	myApi.RunServer()
}
