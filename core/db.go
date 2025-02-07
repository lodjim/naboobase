package core

import (
	"context"
	"errors"
	"fmt"
	"naboobase/configs"
	"naboobase/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (db *MongoDBconnector) DeleteRecordById(ctx context.Context, collectionName string, id primitive.ObjectID, record interface{}) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
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

	fields := utils.GetTaggedFields(record, "autogenerate")
	for _, field := range fields {
		err := utils.Set(field, primitive.NewObjectID(), record)
		if err != nil {
			return err
		}
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

func (db *MongoDBconnector) UpdateRecord(
	ctx context.Context,
	collectionName string,
	id primitive.ObjectID,
	updateData interface{},
	record interface{},
) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)

	if err := isUnique(ctx, collection, updateData, "unique"); err != nil {
		return err
	}

	update := bson.M{"$set": updateData}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		opts,
	).Decode(record)

	return err
}

func (db *MongoDBconnector) BulkCreateRecords(
	ctx context.Context,
	collectionName string,
	records []interface{},
) error {

	collection := db.Client.Database(db.DBName).Collection(collectionName)

	// Check uniqueness for all records (batch version would need optimization)
	for _, record := range records {
		if err := isUnique(ctx, collection, record, "unique"); err != nil {
			return err
		}
	}

	_, err := collection.InsertMany(ctx, records)
	return err
}

func (db *MongoDBconnector) SoftDeleteRecord(
	ctx context.Context,
	collectionName string,
	id primitive.ObjectID,
) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err := collection.UpdateByID(ctx, id, update)
	return err
}
func (db *MongoDBconnector) GetPaginatedRecords(
	ctx context.Context,
	collectionName string,
	filter bson.M,
	page int64,
	limit int64,
	sortField string,
	sortOrder int,
	results interface{},
) (int64, error) {
	collection := db.Client.Database(db.DBName).Collection(collectionName)

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	opts := options.Find().
		SetSkip((page - 1) * limit).
		SetLimit(limit).
		SetSort(bson.D{{sortField, sortOrder}})

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	return total, cursor.All(ctx, results)
}

func (db *MongoDBconnector) ExistsRecord(
	ctx context.Context,
	collectionName string,
	filter bson.M,
) (bool, error) {
	collection := db.Client.Database(db.DBName).Collection(collectionName)
	count, err := collection.CountDocuments(ctx, filter, options.Count().SetLimit(1))
	return count > 0, err
}

func (db *MongoDBconnector) EnsureIndexes(
	ctx context.Context,
	collectionName string,
	model mongo.IndexModel,
) error {
	collection := db.Client.Database(db.DBName).Collection(collectionName)
	_, err := collection.Indexes().CreateOne(ctx, model)
	return err
}

func (db *MongoDBconnector) WithTransaction(
	ctx context.Context,
	fn func(sessCtx mongo.SessionContext) (interface{}, error),
) error {
	session, err := db.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, fn)
	return err
}
