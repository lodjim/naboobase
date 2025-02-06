package controllers

import (
	"github.com/gin-gonic/gin"
	"naboobase/core"
	"naboobase/models"
)

// CreateArticle creates a new Article in the database
func CreateArticle(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateCreateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.ArticleRequest{} },
		NewModel:    func() interface{} { return &models.Article{} },
		NewResponse: func() interface{} { return &models.ArticleResponse{} },
		Collection:  "article",
		Preprocess:  nil,
	})
}

func init() {
	core.AutoEndpointFuncRegistry["article"] = CreateArticle
}
