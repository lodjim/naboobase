package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type VolunterResponse struct {
	Location               string             `json:"location" bson:"location" validate:"max=255"`
	MaritalStatus          string             `json:"marital_status" bson:"marital_status" validate:"max=255"`
	CniRecto               string             `json:"cni_recto" bson:"cni_recto" validate:"max=255"`
	CertificateOfResidence string             `json:"certificate_of_residence" bson:"certificate_of_residence" validate:"max=255"`
	Diplome                string             `json:"diplome" bson:"diplome" validate:"max=255"`
	Id                     primitive.ObjectID `json:"_id" bson:"_id" db:"autogenerate" validate:"max=255"`
	FirstName              string             `json:"first_name" bson:"first_name" validate:"max=255"`
	LastName               string             `json:"last_name" bson:"last_name" validate:"max=255"`
	Cni                    string             `json:"cni" bson:"cni" validate:"max=255"`
	Sex                    string             `json:"sex" bson:"sex" db:"unique" validate:"max=255"`
	BirthDay               string             `json:"birth_day" bson:"birth_day" validate:"max=255"`
	EducationLevel         string             `json:"education_level" bson:"education_level" validate:"max=255"`
	CreatedAt              string             `json:"created_at" bson:"created_at" validate:"max=255"`
	PlaceOfBirth           string             `json:"place_of_birth" bson:"place_of_birth" validate:"max=255"`
	OtherTrainings         string             `json:"other_trainings" bson:"other_trainings" validate:"max=255"`
	CniVerso               string             `json:"cni_verso" bson:"cni_verso" validate:"max=255"`
}
