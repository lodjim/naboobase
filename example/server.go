package main

import (
	"naboobase/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world")
	}
}

func main() {
	myApi := core.Server{}
	myApi.Init("localhost", 1555)
	myApi.AttachEndpoints([]core.Endpoint{
		{
			Method:  "GET",
			Path:    "/health",
			Handler: HealthCheck,
		},
	})
	myApi.RunServer()
}
