package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"naboobase/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForeignKeyConfig struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

type ContentConfig struct {
	ForeignKeys []ForeignKeyConfig `json:"foreign_keys"`
}

type ModelConfig struct {
	ContentConfigs ContentConfig `json:"_config"`
}

type HandlerConfig struct {
	NewRequest  func() interface{}                   // Function to create a new request instance
	NewModel    func() interface{}                   // Function to create a new model instance
	NewResponse func() interface{}                   // Function to create a new response instance
	Collection  string                               // MongoDB collection name
	Preprocess  func(interface{}, interface{}) error // Custom preprocessing function
}

func GenerateCreateHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RequireAuth(c)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Initialize request, model, and response using the provided functions
		req := config.NewRequest()
		model := config.NewModel()
		res := config.NewResponse()
		var modelConfig ModelConfig
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

		jsonData, err := ioutil.ReadFile(fmt.Sprintf("../../json/%s.json", config.Collection))
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(jsonData, &modelConfig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		for _, relations := range modelConfig.ContentConfigs.ForeignKeys {
			if relations.Model == "user" {
				got_claims, ok := c.Get("claims")
				if !ok {
					c.String(http.StatusInternalServerError, "Can't get the ID of the user")
					return
				}
				var claims utils.Claims = got_claims.(utils.Claims)
				utils.Set(relations.Name, claims.Id, &model)
			}
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

func GenerateGetHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		id := c.Param("id")
		req, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		model := config.NewModel()
		res := config.NewModel()
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		err = db.GetRecord(ctx, config.Collection, bson.M{"_id": req}, res)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to get the record: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GenerateGetOneHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		id := c.Param("id")
		req, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		model := config.NewModel()
		res := config.NewModel()
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		err = db.GetRecord(ctx, config.Collection, bson.M{"_id": req}, res)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to get the record: "+err.Error())
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GenerateDeleteHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		id := c.Param("id")
		req, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		model := config.NewModel()
		res := config.NewModel()
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		err = db.DeleteRecordById(ctx, config.Collection, req, res)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to delete the record: "+err.Error())
			return
		}
		utils.Set("_id", id, &res)
		c.String(http.StatusOK, fmt.Sprintf("%s is deleted successfully", id))
	}
}

func GenerateUpdateHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		model := config.NewModel()
		id := c.Param("id")
		req, err := primitive.ObjectIDFromHex(id)
		rawData, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read request body"})
			return
		}
		var data map[string]interface{}
		var modelJson map[string]interface{}
		jsonBytes, err := json.Marshal(model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON"})
			return
		}
		if err := json.Unmarshal(rawData, &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if err = json.Unmarshal(jsonBytes, &modelJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if err := utils.ValidateKeys(data, modelJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		err = db.UpdateRecord(ctx, config.Collection, req, data, model)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to update the record: "+err.Error())
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("%s is updated successfully", id))
	}
}

func GenerateGetAllHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Parse pagination parameters
		page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.ParseInt(c.Query("limit"), 10, 64)
		if limit < 1 {
			limit = 50
		}
		if 10000 > 0 && limit > 10000 {
			limit = 10000
		}

		// Parse sorting parameters
		sortOrder, _ := strconv.Atoi(c.Query("sort_order"))
		if sortOrder != 1 && sortOrder != -1 {
			sortOrder = -1
		}

		sortField := c.Query("sort_field")
		if sortField == "" {
			sortField = "_id"
		}
		/*
			if len(config.AllowedSortFields) > 0 && !contains(config.AllowedSortFields, sortField) {
				c.String(http.StatusBadRequest, "invalid sort field")
				return
				}*/

		// Parse request into req object
		req := config.NewRequest()
		if err := c.ShouldBindQuery(req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		// Initialize model and build filter
		model := config.NewModel()
		var filter bson.M
		filter = bson.M{}
		/*
			if config.BuildFilter != nil {
				filter, err := config.BuildFilter(req, model)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
			} else {
				filter = bson.M{}
			}
		*/
		// Preprocess hook
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		// Retrieve paginated records
		var results []map[string]interface{}
		total, err := db.GetPaginatedRecords(ctx, config.Collection, filter, page, limit, sortField, sortOrder, &results)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total": total,
			"data":  results,
		})
	}
}

// Helper function to check slice containment
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
