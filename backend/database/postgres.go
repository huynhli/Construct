package database

import (
	"backend/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// var Client *mongo.Client
// var ConstructDatabase *mongo.Database
var DB *sql.DB

func ConnectPostgres() {
	if config.DB_PASSWORD == "" {
		log.Fatalf("Empty db creds.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uri := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		config.DB_USERNAME,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)

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
}

func DisconnectPostgres() {
	err := DB.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres connection closed")
}
