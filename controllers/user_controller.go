package controllers

import (
	"naboobase/core"
	"naboobase/models"
	"naboobase/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateUser(db core.MongoDBconnector) gin.HandlerFunc {
	return core.GenerateCreateHandler(db, core.HandlerConfig{
		NewRequest:  func() interface{} { return &models.UserRequest{} },
		NewModel:    func() interface{} { return &models.User{} },
		NewResponse: func() interface{} { return &models.UserResponse{} },
		Collection:  "users",
		Preprocess: func(model, req interface{}) error {
			userReq := req.(*models.UserRequest)
			user := model.(*models.User)
			hashedPassword, err := utils.HashPassword(userReq.Password)
			if err != nil {
				return err
			}
			user.PasswordHashed = hashedPassword
			return nil
		},
	})
}
