package controllers

import (
	"context"
	"fmt"
	"naboobase/core"
	"naboobase/models"
	"naboobase/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	//"reflect"
	"time"
)

func Login(db core.MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var payload models.LoginRequest

		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{Status: http.StatusUnprocessableEntity, ErrorMessage: "Invalid Format Request"})
			return
		}

		if validationErr := validate.Struct(&payload); validationErr != nil {
			c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{Status: http.StatusUnprocessableEntity, ErrorMessage: fmt.Sprintf("Error during the validation of the request: %s", validationErr.Error())})
			return
		}

		var user models.User

		err := db.GetRecord(ctx, "users", bson.M{"email": payload.Email}, &user)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: "User is not found"})
			return
		}
		if err := utils.VerifyPassword(user.PasswordHashed, payload.Password); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Password is not correct"})
			return
		}
		token, err := utils.CreateToken(payload.Email, user.ID.Hex(), user.IsVerified, user.IsSuperuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: "Error while creating access token"})
			return
		}
		c.JSON(http.StatusOK, models.LoginResponse{Token: token, TokenType: "Bearer"})
	}
}
