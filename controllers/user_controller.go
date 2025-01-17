package controllers

import (
	"context"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"naboobase/core"
	"naboobase/models"
	"naboobase/utils"
	"net/http"
	//"reflect"
	"time"
)

var validate = validator.New()

func CreateUser(db core.MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var userRequest models.UserRequest
		var userResponse models.UserResponse
		var user models.User

		if err := c.BindJSON(&userRequest); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if err := validate.Struct(userRequest); err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok {
				c.String(http.StatusBadRequest, "Invalid request payload")
				return
			}
			c.String(http.StatusBadRequest, validationErrors.Error())
			return
		}

		if err := copier.Copy(&user, &userRequest); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		hashedPassword, err := utils.HashPassword(userRequest.Password)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		user.PasswordHashed = hashedPassword

		if err := db.CreateRecord(ctx, "user", user); err != nil {
			c.String(http.StatusBadRequest, "Failed to create user: "+err.Error())
			return
		}

		if err := copier.Copy(&userResponse, &user); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, userResponse)
	}
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
