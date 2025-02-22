package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/core"
	"github.com/lodjim/naboobase/models"
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

func GetOneArticle(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetOneHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.ArticleRequest{} },
		NewModel:    func() interface{} { return &models.Article{} },
		NewResponse: func() interface{} { return &models.ArticleResponse{} },
		Collection:  "article",
		Preprocess:  nil,
	})
}
func GetAllArticle(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetAllHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.ArticleRequest{} },
		NewModel:    func() interface{} { return &models.Article{} },
		NewResponse: func() interface{} { return &models.ArticleResponse{} },
		Collection:  "article",
		Preprocess:  nil,
	})
}

func UpdateArticle(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateUpdateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.ArticleRequest{} },
		NewModel:    func() interface{} { return &models.Article{} },
		NewResponse: func() interface{} { return &models.ArticleResponse{} },
		Collection:  "article",
		Preprocess:  nil,
	})
}

func DeleteArticle(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateDeleteHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.ArticleRequest{} },
		NewModel:    func() interface{} { return &models.Article{} },
		NewResponse: func() interface{} { return &models.ArticleResponse{} },
		Collection:  "article",
		Preprocess:  nil,
	})
}

func init() {
	core.AutoEndpointFuncRegistry["article-POST"] = CreateArticle
	core.AutoEndpointFuncRegistry["article-GET-ID"] = GetOneArticle
	core.AutoEndpointFuncRegistry["article-GET"] = GetAllArticle
	core.AutoEndpointFuncRegistry["article-PUT-ID"] = UpdateArticle
	core.AutoEndpointFuncRegistry["article-DELETE-ID"] = DeleteArticle
}
