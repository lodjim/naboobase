package models

type UserRequest struct {
	Password string `json:"password" bson:"password" validate:"required,min=8,max=16"`
	Name     string `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
}
