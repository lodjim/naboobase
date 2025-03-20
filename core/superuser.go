package core

import (
	"github.com/lodjim/naboobase/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"github.com/lodjim/naboobase/models"
	//"reflect"
)

type SuperUserManagement struct {
}

func (superUserManagement *SuperUserManagement) CreateSuperUser(db MongoDBconnector) gin.HandlerFunc {
	return GenerateCreateHandler(db, HandlerConfig{
		NewRequest:    func() interface{} { return &models.UserRequest{} },
		NewModel:      func() interface{} { return &models.User{} },
		NewResponse:   func() interface{} { return &models.UserResponse{} },
		Functionality: "superuser",
		Collection:    "user",
		Preprocess: func(model, req interface{}, filter *bson.M) error {
			userReq := req.(*models.UserRequest)
			user := model.(*models.User)
			hashedPassword, err := utils.HashPassword(userReq.Password)
			if err != nil {
				return err
			}
			user.PasswordHashed = hashedPassword
			user.IsSuperuser = true
			return nil
		},
	})
}

func (superUserManagement *SuperUserManagement) GetsSuperUsers(db MongoDBconnector) gin.HandlerFunc {
	return GenerateGetAllHandler(db, HandlerConfig{
		NewRequest:    func() interface{} { return &models.UserRequest{} },
		NewModel:      func() interface{} { return &models.User{} },
		NewResponse:   func() interface{} { return &models.UserResponse{} },
		Functionality: "superuser",
		Collection:    "user",
		Preprocess: func(model, req interface{}, filter *bson.M) error {
			if filter == nil {
				emptyFilter := bson.M{}
				filter = &emptyFilter
			}
			(*filter)["is_superuser"] = true
			return nil
		},
	})
}

func (superUserManagement *SuperUserManagement) Init(db MongoDBconnector) []Endpoint {
	var endpoints []Endpoint = []Endpoint{{
		Method:  "POST",
		Path:    "/admin/superuser",
		Handler: superUserManagement.CreateSuperUser(db),
	},
		{
			Method:  "GET",
			Path:    "/admin/superuser",
			Handler: superUserManagement.GetsSuperUsers(db),
		},
	}
	return endpoints
}
