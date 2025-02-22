package models

type ArticleResponse struct {
	Description string `json:"description" bson:"description" validate:"max=255"`
	ID          string `json:"_id" bson:"_id" validate:"max=255"`
	Title       string `json:"title" bson:"title" validate:"max=255"`
}
