package models

type UserResponse struct {
	Id        string `json:"_id" bson:"_id" validate:"uuid"`
	Name      string `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
	Role      string `json:"role" bson:"role" validate:"required"`
	CreatedAt string `json:"created_at" bson:"created_at" validate:"required"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}
