package controllers

import (
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"naboobase/core"
	"naboobase/models"
	"naboobase/utils"
	//"reflect"
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

/*
func CreateUser(db core.MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var userRequest models.UserRequest
		var userResponse models.UserResponse
		var user models.User
		err := c.BindJSON(&userResponse)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

			    userType, err := utils.GetModelType("user")
					if err != nil {
						fmt.Println(err)
						c.String(http.StatusBadRequest, err.Error())
						return
					}
					user := reflect.New(userType).Interface()


		err = validate.Struct(userRequest)
		if err != nil {
			validationErrors := err.(validator.ValidationErrors)
			c.String(http.StatusBadRequest, validationErrors.Error())
			return
		}

		user.PasswordHashed, err = utils.HashPassword(userRequest.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		err = copier.Copy(&user, userRequest)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		err = db.CreateRecord(ctx, "user", user)
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		err = copier.Copy(&userResponse, user)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, userResponse)
	}
}*/
