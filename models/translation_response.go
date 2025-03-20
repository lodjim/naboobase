package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TranslationResponse struct {
	Wolof     string             `json:"wolof" bson:"wolof" validate:"max=255"`
	French    string             `json:"french" bson:"french" validate:"max=255"`
	IsGood    string             `json:"is_good" bson:"is_good" validate:"max=255"`
	CreatedAt string             `json:"created_at" bson:"created_at" db:"none" validate:"max=255"`
	Id        primitive.ObjectID `json:"_id" bson:"_id" db:"autogenerate" validate:"max=255"`
}
