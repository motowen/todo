package database

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var todoCollection *mongo.Collection

var ERROR_DATA_NOT_FOUND = errors.New("data not found")

func Setup(uri string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		return
	}

	databaseName := cs.Database
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		return
	}

	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		return
	}

	todoCollection = client.Database(databaseName).Collection("validation")

	return
}

func Drop() (err error) {
	if err = todoCollection.Drop(context.TODO()); err != nil {
		return
	}
	return
}
