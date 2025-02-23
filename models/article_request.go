package models

type ArticleRequest struct {
	Title       string `json:"title" bson:"title" validate:"max=255"`
	Description string `json:"description" bson:"description" validate:"max=255"`
}
