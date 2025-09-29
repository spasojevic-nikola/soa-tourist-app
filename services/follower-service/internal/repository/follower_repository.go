package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)


type FollowerRepository struct {
	Driver neo4j.DriverWithContext
}

// NewFollowerRepository kreira novu instancu repozitorijuma
func NewFollowerRepository(driver neo4j.DriverWithContext) *FollowerRepository {
	return &FollowerRepository{Driver: driver}
}

// Follow kreira :FOLLOWS vezu između dva korisnika
func (repo *FollowerRepository) Follow(followerId, followedId uint) error {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)
	//Cypher upit
	query := `
		MERGE (follower:User {id: $followerId})
		MERGE (followed:User {id: $followedId})
		MERGE (follower)-[:FOLLOWS]->(followed)
	`
	params := map[string]interface{}{
		"followerId": followerId,
		"followedId": followedId,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}
