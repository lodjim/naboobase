package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `json:"_id" bson:"_id" db:"autogenerate" validate:"max=255"`
	Name           string             `json:"name" bson:"name" validate:"max=255"`
	Email          string             `json:"email" bson:"email" db:"unique" validate:"email"`
	IsVerified     bool               `json:"is_verified" bson:"is_verified"`
	CreatedAt      string             `json:"created_at" bson:"created_at" validate:"max=255"`
	PasswordHashed string             `json:"passwordHashed" bson:"passwordHashed" validate:"max=255"`
	IsSuperuser    bool               `json:"is_superuser" bson:"is_superuser"`
	Role           string             `json:"role" bson:"role" validate:"max=255"`
	UpdatedAt      string             `json:"updated_at" bson:"updated_at" validate:"max=255"`
}
