package core

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lodjim/naboobase/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lodjim/naboobase/models"
	"go.mongodb.org/mongo-driver/bson"

	//"reflect"
	"time"
)

var validate = validator.New()

type AuthenticatorType string

const (
	PasswordTypeAuthenticator AuthenticatorType = "password"
)

type Authenticator struct {
	Type AuthenticatorType
}

func (Auth *Authenticator) Init(db MongoDBconnector) []Endpoint {
	var endpoints []Endpoint = []Endpoint{{
		Method:  "POST",
		Path:    "/login",
		Handler: Auth.Login(db),
	},
		{
			Method:  "POST",
			Path:    "/refresh-token",
			Handler: Auth.RefreshToken(db),
		},
	}
	return endpoints
}

func (Auth *Authenticator) Login(db MongoDBconnector) gin.HandlerFunc {
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

		err := db.GetRecord(ctx, "user", bson.M{"email": payload.Email}, &user)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: "User is not found"})
			return
		}
		if err := utils.VerifyPassword(user.PasswordHashed, payload.Password); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Status: http.StatusBadRequest, ErrorMessage: "Password is not correct"})
			return
		}
		token, err := utils.CreateToken(payload.Email, user.ID.Hex(), user.IsVerified, user.IsSuperuser)
		refreshToken, err := utils.CreateRefreshToken(payload.Email, user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: "Error while creating access token"})
			return
		}
		c.JSON(http.StatusOK, models.LoginResponse{Token: token, TokenType: "Bearer", RefreshToken: refreshToken})
	}
}

func (Auth *Authenticator) RefreshToken(db MongoDBconnector) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var payload models.RefreshTokenRequest
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{Status: http.StatusUnprocessableEntity, ErrorMessage: "Invalid Format Request"})
			return
		}
		if validationErr := validate.Struct(&payload); validationErr != nil {
			c.JSON(http.StatusUnprocessableEntity, models.ErrorResponse{Status: http.StatusUnprocessableEntity, ErrorMessage: fmt.Sprintf("Error during the validation of the request: %s", validationErr.Error())})
			return
		}
		claims, err := utils.VerifyRefreshJWT(payload.RefreshToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Status: http.StatusUnauthorized, ErrorMessage: err.Error()})
			return
		}
		var user models.User
		err = db.GetRecord(ctx, "user", bson.M{"email": claims.Email}, &user)
		if err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Status: http.StatusNotFound, ErrorMessage: "User is not found"})
			return
		}
		token, err := utils.CreateToken(claims.Email, user.ID.Hex(), user.IsVerified, user.IsSuperuser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: "Error while creating access token"})
			return
		}
		c.JSON(http.StatusOK, models.LoginResponse{Token: token, TokenType: "Bearer", RefreshToken: payload.RefreshToken})
	}
}
