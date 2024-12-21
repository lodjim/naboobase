package core

import "github.com/gin-gonic/gin"

type Endpoint struct {
	Method  string
	Path    string
	Handler *gin.HandlerFunc
}
