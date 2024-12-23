package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"naboobase/proto_struct"
)

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user proto_struct.User
		defer cancel()
		err := c.BindJSON(&user)

		/*
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var user requests.UserRequest
			defer cancel()

			if err := c.BindJSON(&user); err != nil {
				c.JSON(http.StatusUnprocessableEntity, responses.ErrorResponse{Status: http.StatusUnprocessableEntity, ErrorMessage: fmt.Sprintf("The form is not valid detail: s% ", err.Error())})
				return
			}

			if validationErr := validate.Struct(&user); validationErr != nil {
				c.JSON(http.StatusBadRequest, responses.ErrorResponse{Status: http.StatusBadRequest,
					ErrorMessage: fmt.Sprintf("Error during the validation of the form detail: %s", validationErr.Error())})
				return
			}
			id_user := primitive.NewObjectID()
			hashed_password, _ := utils.HashPassword(user.Password)
			newUser := models.User{
				Id:            id_user,
				Email:         user.Email,
				LastName:      user.LastName,
				FirstName:     user.FirstName,
				Role:          constants.Owner,
				PhoneNumber:   user.PhoneNumber,
				Password:      hashed_password,
				Enable2Fac:    false,
				IsActivated:   false,
				Address:       user.Address,
				Country:       string(user.Country),
				CreatedAt:     time.Now().Format(time.RFC3339),
				UpdatedAt:     time.Now().Format(time.RFC3339),
				LastLogin:     time.Now().Format(time.RFC3339),
				Organizations: make([]primitive.ObjectID, 0),
			}
			_, err := userCollection.InsertOne(ctx, newUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: fmt.Sprintf("Error during database calling detail: %s", err.Error())})
				return
			}

			expiryDate := time.Now().Add(10 * 24 * time.Hour)

			newCreditBalance := models.CreditBalance{
				ID:         primitive.NewObjectID(),
				UserID:     id_user,
				PackName:   "Pack50K",
				Credits:    0,
				ExpiryDate: expiryDate.Format(time.RFC3339),
			}
			_, err = creditBalanceCollection.InsertOne(ctx, newCreditBalance)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Status: http.StatusInternalServerError, ErrorMessage: fmt.Sprintf("Error during database calling detail: %s", err.Error())})
				return
			}
		*/
		//c.JSON(http.StatusCreated, responses.UserCreationResponse{Id: id_user, Email: user.Email})
		c.String(http.StatusOK, "test")
	}
}
