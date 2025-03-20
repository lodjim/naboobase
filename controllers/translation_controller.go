package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/core"
	"github.com/lodjim/naboobase/models"
)

// CreateTranslation creates a new Translation in the database
func CreateTranslation(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateCreateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.TranslationRequest{} },
		NewModel:    func() interface{} { return &models.Translation{} },
		NewResponse: func() interface{} { return &models.TranslationResponse{} },
		Functionality: "translation",
		Collection:  "translation",
		Preprocess:  nil,
	})
}

func GetOneTranslation(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetOneHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.TranslationRequest{} },
		NewModel:    func() interface{} { return &models.Translation{} },
		NewResponse: func() interface{} { return &models.TranslationResponse{} },
		Functionality: "translation",
		Collection:  "translation",
		Preprocess:  nil,
	})
}
func GetAllTranslation(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetAllHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.TranslationRequest{} },
		NewModel:    func() interface{} { return &models.Translation{} },
		NewResponse: func() interface{} { return &models.TranslationResponse{} },
		Functionality: "translation",
		Collection:  "translation",
		Preprocess:  nil,
	})
}

func UpdateTranslation(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateUpdateHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.TranslationRequest{} },
			NewModel:    func() interface{} { return &models.Translation{} },
			NewResponse: func() interface{} { return &models.TranslationResponse{} },
			Functionality: "translation",
			Collection:  "translation",
			Preprocess:  nil,
	})
}

func DeleteTranslation(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateDeleteHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.TranslationRequest{} },
			NewModel:    func() interface{} { return &models.Translation{} },
			NewResponse: func() interface{} { return &models.TranslationResponse{} },
			Functionality: "translation",
			Collection:  "translation",
			Preprocess:  nil,
	})
}


func init() {
	core.AutoEndpointFuncRegistry["translation-POST"] = CreateTranslation
	core.AutoEndpointFuncRegistry["translation-GET-ID"] = GetOneTranslation
	core.AutoEndpointFuncRegistry["translation-GET"] = GetAllTranslation
	core.AutoEndpointFuncRegistry["translation-PUT-ID"] = UpdateTranslation
	core.AutoEndpointFuncRegistry["translation-DELETE-ID"] = DeleteTranslation
}
