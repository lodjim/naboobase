package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/lodjim/naboobase/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ForeignKeyConfig struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

type AuthRules struct {
	ShouldBeAuthenticated bool `json:"should_be_authenticated"`
	OnlyForAdmin          bool `json:"only_for_admin"`
}

type CRUDConfig struct {
	AuthRules AuthRules `json:"auth_rules"`
}

type ContentConfig struct {
	ForeignKeys []ForeignKeyConfig `json:"foreign_keys"`
	Create      CRUDConfig         `json:"create"`
	Delete      CRUDConfig         `json:"delete"`
	Update      CRUDConfig         `json:"update"`
	GetOne      CRUDConfig         `json:"getOne"`
	GetAll      CRUDConfig         `json:"getAll"`
}

type ModelConfig struct {
	ContentConfigs ContentConfig `json:"_config"`
}

type HandlerConfig struct {
	NewRequest    func() interface{} // Function to create a new request instance
	NewModel      func() interface{} // Function to create a new model instance
	NewResponse   func() interface{} // Function to create a new response instance
	Functionality string
	Collection    string                                        // MongoDB collection name
	Preprocess    func(interface{}, interface{}, *bson.M) error // Custom preprocessing function
}

func loadConfig(collectionName string, modelConfig *ModelConfig) error {
	jsonData, err := ioutil.ReadFile(fmt.Sprintf("./json/%s.json", collectionName))
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}
	if err := json.Unmarshal(jsonData, modelConfig); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}
	return nil
}

func GenerateCreateHandler(db MongoDBconnector, config HandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var claims *utils.Claims
		var modelConfig ModelConfig
		if err := loadConfig(config.Collection, &modelConfig); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if modelConfig.ContentConfigs.Create.AuthRules.ShouldBeAuthenticated {
			if config.Functionality != "user" {
				utils.RequireAuth(c)
			}
		}

		got_claims, ok := c.Get("claims")
		if modelConfig.ContentConfigs.Create.AuthRules.OnlyForAdmin {
			if !ok {
				c.String(http.StatusInternalServerError, "Can't get the ID of the user")
				return
			}
			claims = got_claims.(*utils.Claims)
			if !claims.IsSuperUser {
				c.String(http.StatusUnauthorized, "You are not admin")
				return
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := config.NewRequest()
		model := config.NewModel()
		res := config.NewResponse()

		if err := c.BindJSON(req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if err := validate.Struct(req); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				c.String(http.StatusBadRequest, validationErrors.Error())
			} else {
				c.String(http.StatusBadRequest, "Invalid request payload")
			}
			return
		}

		if config.Collection != "user" {
			for _, relations := range modelConfig.ContentConfigs.ForeignKeys {
				if relations.Model == "user" {
					if !ok {
						c.String(http.StatusInternalServerError, "Can't get the ID of the user")
						return
					}
					claims = got_claims.(*utils.Claims)
					utils.Set(relations.Name, claims.Id, &model)
				}
			}
		}
		// Copy data from request to model using copier
		if err := copier.Copy(model, req); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Execute custom preprocessing (e.g., password hashing)
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req, nil); err != nil {
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
			if err := config.Preprocess(model, req, nil); err != nil {
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
		var claims *utils.Claims
		var modelConfig ModelConfig
		if err := loadConfig(config.Collection, &modelConfig); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if modelConfig.ContentConfigs.GetOne.AuthRules.ShouldBeAuthenticated {
			if config.Functionality != "user" {
				utils.RequireAuth(c)
			}
		}
		got_claims, ok := c.Get("claims")
		if modelConfig.ContentConfigs.GetOne.AuthRules.OnlyForAdmin {
			if !ok {
				c.String(http.StatusInternalServerError, "Can't get the ID of the user")
				return
			}
			claims = got_claims.(*utils.Claims)
			if !claims.IsSuperUser {
				c.String(http.StatusUnauthorized, "You are not admin")
				return
			}
		}
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
			if err := config.Preprocess(model, req, nil); err != nil {
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
		var claims *utils.Claims
		var modelConfig ModelConfig
		if err := loadConfig(config.Collection, &modelConfig); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if modelConfig.ContentConfigs.Delete.AuthRules.ShouldBeAuthenticated {
			if config.Functionality != "user" {
				utils.RequireAuth(c)
			}
		}
		got_claims, ok := c.Get("claims")
		if modelConfig.ContentConfigs.Delete.AuthRules.OnlyForAdmin {
			if !ok {
				c.String(http.StatusInternalServerError, "Can't get the ID of the user")
				return
			}
			claims = got_claims.(*utils.Claims)
			if !claims.IsSuperUser {
				c.String(http.StatusUnauthorized, "You are not admin")
				return
			}
		}

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
			if err := config.Preprocess(model, req, nil); err != nil {
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
		var claims *utils.Claims
		var modelConfig ModelConfig
		if err := loadConfig(config.Collection, &modelConfig); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if modelConfig.ContentConfigs.Update.AuthRules.ShouldBeAuthenticated {
			if config.Functionality != "user" {
				utils.RequireAuth(c)
			}
		}
		got_claims, ok := c.Get("claims")
		if modelConfig.ContentConfigs.Update.AuthRules.OnlyForAdmin {
			if !ok {
				c.String(http.StatusInternalServerError, "Can't get the ID of the user")
				return
			}
			claims = got_claims.(*utils.Claims)
			if !claims.IsSuperUser {
				c.String(http.StatusUnauthorized, "You are not admin")
				return
			}
		}

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
			if err := config.Preprocess(model, req, nil); err != nil {
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
		var claims *utils.Claims
		var modelConfig ModelConfig
		if err := loadConfig(config.Collection, &modelConfig); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if modelConfig.ContentConfigs.GetAll.AuthRules.ShouldBeAuthenticated {
			if config.Functionality != "user" {
				utils.RequireAuth(c)
			}
		}
		got_claims, ok := c.Get("claims")
		fmt.Println(got_claims)
		if modelConfig.ContentConfigs.GetAll.AuthRules.OnlyForAdmin {
			if !ok {
				c.String(http.StatusInternalServerError, "Can't get the ID of the user")
				return
			}
			claims = got_claims.(*utils.Claims)
			if !claims.IsSuperUser {
				c.String(http.StatusUnauthorized, "You are not admin")
				return
			}
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		filter_search := c.Query("filter")
		fmt.Println(filter_search)

		var filter *bson.M
		if filter_search != "" {
			query, err := utils.TransformFilterToMongoQuery(filter_search)
			if err != nil {
				c.String(http.StatusBadRequest, "The filter used is not appropriate")
				return
			}
			filter = &query
		} else {
			emptyQuery := bson.M{}
			filter = &emptyQuery
		}

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

		sortOrder, _ := strconv.Atoi(c.Query("sort_order"))
		if sortOrder != 1 && sortOrder != -1 {
			sortOrder = -1
		}

		sortField := c.Query("sort_field")
		if sortField == "" {
			sortField = "_id"
		}

		req := config.NewRequest()
		if err := c.ShouldBindQuery(req); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		model := config.NewModel()
		if config.Preprocess != nil {
			if err := config.Preprocess(model, req, filter); err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}
		response := config.NewResponse()
		var results []map[string]interface{}
		total, err := db.GetPaginatedRecords(ctx, config.Collection, *filter, page, limit, sortField, sortOrder, &results)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		var resultToReturn []map[string]interface{}

		for _, item := range results {
			var newItem map[string]interface{}
			for key, value := range item {
				value, err = utils.Get(key, response)
				fmt.Println("some thing--------------------")
				fmt.Println(value)
				if err != nil {
					continue
				}
				newItem[key] = value
			}
			resultToReturn = append(resultToReturn, newItem)
		}
		c.JSON(http.StatusOK, gin.H{
			"total": total,
			"data":  resultToReturn,
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
