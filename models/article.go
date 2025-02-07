package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Article struct {
	CreatedAt   string             `json:"created_at" bson:"created_at" validate:"max=255"`
	UpdatedAt   string             `json:"updated_at" bson:"updated_at" validate:"max=255"`
	ID          primitive.ObjectID `json:"_id" bson:"_id" db:"autogenerate" validate:"max=255"`
	Title       string             `json:"title" bson:"title" validate:"max=255"`
	Description string             `json:"description" bson:"description" validate:"max=255"`
}
