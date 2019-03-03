package db

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type collections struct {
	Comics *mongo.Collection
}

// DB is a Database
var DB *mongo.Database

// Collections defines the collections can be used.
var Collections collections

// Connect to Database
func Connect(mongoURL string, dbName string) (*mongo.Client, error) {
	client, e := mongo.Connect(context.Background(), mongoURL, nil)

	if e != nil {
		return nil, e
	}

	DB = client.Database(dbName)
	Collections.Comics = DB.Collection("comics")

	return client, nil
}
