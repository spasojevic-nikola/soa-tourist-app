package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDB uspostavlja konekciju sa MongoDB bazom.
func InitDB() *mongo.Database {
	// Koristi ime 'mongo' servisa iz docker-compose.yml
	dsn := "mongodb://mongo:27017" 

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ Failed to ping MongoDB. DB is unavailable: %v", err)
	}

	log.Println("✅ Successfully connected to MongoDB!")

	// Vraćamo bazu 'purchase_db'
	return client.Database("purchase_db")
}