package models

type UserRequest struct {
	Name     string `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required,min=8,max=16"`
}
