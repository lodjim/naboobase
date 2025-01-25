package core

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"naboobase/configs"
	"naboobase/utils"
	"time"
)

type DBconnector interface {
	Connect(url string) error
	//Close() error
	CreateRecord() any
	//GetRecord() error
	//DeleteRecord() any
	//UpdateRecord() any
}

type MongoDBconnector struct {
	DBName string
	Client *mongo.Client
}

func isFieldUnique(ctx context.Context, collection *mongo.Collection, field string, value interface{}) (bool, error) {
	// Create a filter to find documents with the specified field value
	filter := bson.M{field: value}
	// Query the collection
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	// If no document is found, the field is unique
	if err == mongo.ErrNoDocuments {
		return true, nil
	} else if err != nil {
		return false, err
	}
	// If a document is found, the field is not unique
	return false, nil
}

func isUnique(ctx context.Context, collection *mongo.Collection, record interface{}, tag string) error {
	fields := utils.GetTaggedFields(record, tag)
	if len(fields) != 0 {
		for _, field := range fields {
			fmt.Println(field)
			value, err := utils.Get(field, record)
			if err != nil {
				return err
			}
			fmt.Println(value)
			ok, err := isFieldUnique(ctx, collection, utils.ConvertToSnakeCase(field), value)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("For the Field: %s the value %s is already in the database", field, value)
			}
			fmt.Println("the value of isFieldUnique ")
			fmt.Println(ok)
		}
	}
	return nil
}

func (db *MongoDBconnector) GetRecord(ctx context.Context, collectionName string, filter interface{}, record interface{}) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(record)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDBconnector) CreateRecord(ctx context.Context, collectionName string, record interface{}) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)

	err := isUnique(ctx, collection, record, "unique")

	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, record)
	if err != nil {
		return err
	}
	return nil
}

func (db *MongoDBconnector) Connect(DBName string) error {
	db.DBName = DBName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	db.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.EnvMongoURI()).SetMaxPoolSize(50).SetMinPoolSize(10).SetMaxConnIdleTime(10*time.Minute))
	if err != nil {
		return err
	}
	err = db.Client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Connected to MongoDB")
	return errors.New("Error during the database connection")
}
