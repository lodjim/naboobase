package core

import (
	"github.com/gin-gonic/gin"
)

var AutoEndpointFuncRegistry = make(map[string]func(MongoDBconnector) gin.HandlerFunc)
