package main

import (
	"naboobase/controllers"
	"naboobase/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
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
			Path:    "/test",
			Handler: controllers.CreateUser(dbConnector),
		},
	})
	myApi.RunServer()
}
