package database

import (
	"backend/config"
	"context"
	"database/sql"
	"log"
	"time"
)

// var Client *mongo.Client
// var ConstructDatabase *mongo.Database
var DB *sql.DB

func ConnectPostgres() {
	if config.DB_PASSWORD == "" {
		log.Fatalf("Empty db creds.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "postgresql://postgres:" + config.DB_PASSWORD + "@db.iahttzkfroabvsbofbpr.supabase.co:5432/postgres"
	var err error
	DB, err = sql.Open("pgx", uri)
	if err != nil {
		log.Fatalf("Error opening Postgres connection: %s", err)
	}

	// Ping to verify connection
	err = DB.PingContext(ctx)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %s", err)
	}

	log.Println("Connected to Postgres successfully")

	// mongoDBConfig := options.Client().ApplyURI(uri)
	// var err error
	// Client, err = mongo.Connect(ctx, mongoDBConfig)
	// if err != nil {
	// 	log.Fatalf("Error connecting to MongoDB: %s", err)
	// }
	// ConstructDatabase = Client.Database("Prototype")
}

func DisconnectPostgres() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres connection closed")
	// err := Client.Disconnect(context.TODO())
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
