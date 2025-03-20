package models

type TranslationRequest struct {
	Wolof  string `json:"wolof" bson:"wolof" validate:"max=255"`
	French string `json:"french" bson:"french" validate:"max=255"`
}
