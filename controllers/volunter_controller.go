package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/core"
	"github.com/lodjim/naboobase/models"
)

// CreateVolunter creates a new Volunter in the database
func CreateVolunter(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateCreateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.VolunterRequest{} },
		NewModel:    func() interface{} { return &models.Volunter{} },
		NewResponse: func() interface{} { return &models.VolunterResponse{} },
		Functionality: "volunter",
		Collection:  "volunter",
		Preprocess:  nil,
	})
}

func GetOneVolunter(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetOneHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.VolunterRequest{} },
		NewModel:    func() interface{} { return &models.Volunter{} },
		NewResponse: func() interface{} { return &models.VolunterResponse{} },
		Functionality: "volunter",
		Collection:  "volunter",
		Preprocess:  nil,
	})
}
func GetAllVolunter(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateGetAllHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.VolunterRequest{} },
		NewModel:    func() interface{} { return &models.Volunter{} },
		NewResponse: func() interface{} { return &models.VolunterResponse{} },
		Functionality: "volunter",
		Collection:  "volunter",
		Preprocess:  nil,
	})
}

func UpdateVolunter(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateUpdateHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.VolunterRequest{} },
			NewModel:    func() interface{} { return &models.Volunter{} },
			NewResponse: func() interface{} { return &models.VolunterResponse{} },
			Functionality: "volunter",
			Collection:  "volunter",
			Preprocess:  nil,
	})
}

func DeleteVolunter(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateDeleteHandler(db, core.HandlerConfig{
			NewRequest:  func() interface{} { return &models.VolunterRequest{} },
			NewModel:    func() interface{} { return &models.Volunter{} },
			NewResponse: func() interface{} { return &models.VolunterResponse{} },
			Functionality: "volunter",
			Collection:  "volunter",
			Preprocess:  nil,
	})
}


func init() {
	core.AutoEndpointFuncRegistry["volunter-POST"] = CreateVolunter
	core.AutoEndpointFuncRegistry["volunter-GET-ID"] = GetOneVolunter
	core.AutoEndpointFuncRegistry["volunter-GET"] = GetAllVolunter
	core.AutoEndpointFuncRegistry["volunter-PUT-ID"] = UpdateVolunter
	core.AutoEndpointFuncRegistry["volunter-DELETE-ID"] = DeleteVolunter
}
