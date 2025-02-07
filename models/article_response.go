package models

type ArticleResponse struct {
	ID          string `json:"_id" bson:"_id" validate:"max=255"`
	Title       string `json:"title" bson:"title" validate:"max=255"`
	Description string `json:"description" bson:"description" validate:"max=255"`
}
