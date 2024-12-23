package core

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"naboobase/configs"
	pb "naboobase/proto_struct"
	"reflect"
	"strings"
	"time"
)

func LoadProtoStruct(structName string) (interface{}, error) {
	// Capitalize first letter to match struct name
	structName = strings.Title(strings.ToLower(structName))

	// Get the type by name from the proto package
	t := reflect.ValueOf(pb.File{}).Elem().Type()
	structType := reflect.TypeOf(pb.File{}).PkgPath()

	// Create new instance of the struct
	v := reflect.New(t).Interface()
	return v, nil
}

type DBconnector interface {
	Connect(url string) error
	Close() error
	CreateRecord() any
	GetRecord() any
	DeleteRecord() any
	UpdateRecord() any
}

type MongoDBconnector struct {
	DBName string
	Client *mongo.Client
}

func (db *MongoDBconnector) Connect(DBName string) error {
	db.DBName = DBName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	db.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(configs.EnvMongoURI()))
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
