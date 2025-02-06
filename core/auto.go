package core

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"net/http"
	"time"
)

type HandlerConfig struct {
	NewRequest  func() interface{}                   // Function to create a new request instance
	NewModel    func() interface{}                   // Function to create a new model instance
	NewResponse func() interface{}                   // Function to create a new response instance
	Collection  string                               // MongoDB collection name
	Preprocess  func(interface{}, interface{}) error // Custom preprocessing function
}

func GenerateCreateHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Initialize request, model, and response using the provided functions
		req := config.NewRequest()
		model := config.NewModel()
		res := config.NewResponse()

		// Bind incoming JSON to the request struct
		if err := c.BindJSON(req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		// Validate the request struct using the validator
		if err := validate.Struct(req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.String(http.StatusBadRequest, validationErrors.Error())
			} else {
				c.String(http.StatusBadRequest, "Invalid request payload")
			}
			return
		}

		// Copy data from request to model using copier
		if err := copier.Copy(model, req); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Execute custom preprocessing (e.g., password hashing)
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Insert the model into the database
		if err := db.CreateRecord(ctx, config.Collection, model); err != nil {
			c.String(http.StatusBadRequest, "Failed to create record: "+err.Error())
			return
		}

		// Copy data from model to response
		if err := copier.Copy(res, model); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Return the response as JSON
		c.JSON(http.StatusOK, res)
	}
}
