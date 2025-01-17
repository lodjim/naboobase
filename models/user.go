package models

type User struct {
	Role           string `json:"role" bson:"role" validate:"max=255"`
	CreatedAt      string `json:"created_at" bson:"created_at" validate:"max=255"`
	UpdatedAt      string `json:"updated_at" bson:"updated_at" validate:"max=255"`
	Name           string `json:"name" bson:"name" validate:"max=255"`
	Email          string `json:"email" bson:"email" db:"unique" validate:"email"`
	PasswordHashed string `json:"passwordHashed" bson:"passwordHashed" validate:"max=255"`
}
