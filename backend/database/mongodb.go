package database

import (
	"backend/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func ConnectMongoDB() {
	if config.DB_USERNAME == "" || config.DB_PASSWORD == "" {
		log.Fatalf("Empty db creds.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb+srv://" + config.DB_USERNAME + ":" + config.DB_PASSWORD + "@prototype.iysqpyf.mongodb.net/?retryWrites=true&w=majority&appName=Prototype"
	mongoDBConfig := options.Client().ApplyURI(uri)
	var err error
	Client, err = mongo.Connect(ctx, mongoDBConfig)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err)
	}
	Database = Client.Database("ConstructDB")
}

func DisconnectMongoDB() {
	err := Client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
