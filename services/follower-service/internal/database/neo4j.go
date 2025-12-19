package database

import (
	"context" 	
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// InitDB uspostavlja konekciju sa Neo4j bazom i vraca drajver
func InitDB() neo4j.DriverWithContext {
	uri := os.Getenv("NEO4J_URI")
	user := os.Getenv("NEO4J_USER")
	pass := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		log.Fatalf("FATAL: Could not create Neo4j driver: %s", err)
	}

	// Proveravamo konekciju
	err = driver.VerifyConnectivity(context.Background()) 
	if err != nil {
		log.Fatalf("FATAL: Could not connect to Neo4j: %s", err)
	}

	log.Println("Successfully connected to Neo4j!")
	return driver
}