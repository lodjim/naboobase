package models

import (
	"time"
)

type User struct {
	ID        int         `json:"id" bson:"id" validate:"required,gte=1"`
	Email     string      `json:"email" bson:"email" validate:"required,email"`
	Username  string      `json:"username" bson:"username" validate:"required,min=3,max=30"`
	Age       int         `json:"age" bson:"age" validate:"gte=18,lte=100"`
	Website   string      `json:"website" bson:"website" validate:"required,url"`
	CreatedAt time.Time   `json:"created_at" bson:"created_at" validate:"required,datetime"`
	Profile   UserProfile `json:"profile" bson:"profile" validate:"required"`
	Tags      []string    `json:"tags" bson:"tags" validate:"required,dive,min=2,max=20"`
}

type UserProfile struct {
	Bio      string `json:"bio" bson:"bio" validate:"max=500"`
	Location string `json:"location" bson:"location" validate:"max=100"`
	IsPublic bool   `json:"is_public" bson:"is_public" validate:"required"`
}
