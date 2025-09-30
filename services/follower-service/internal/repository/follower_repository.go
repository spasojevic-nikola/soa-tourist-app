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

// Unfollow brise :FOLLOWS vezu između dva korisnika
func (repo *FollowerRepository) Unfollow(followerId, followedId uint) error {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MATCH (follower:User {id: $followerId})-[r:FOLLOWS]->(followed:User {id: $followedId})
		DELETE r
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
	//proveri jel zapratio
func (repo *FollowerRepository) CheckFollows(followerId, followedId uint) (bool, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		RETURN EXISTS( (:User {id: $followerId})-[:FOLLOWS]->(:User {id: $followedId}) )
	`
	params := map[string]interface{}{
		"followerId": followerId,
		"followedId": followedId,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return false, err
		}
		if res.Next(ctx) {
			return res.Record().Values[0], nil
		}
		return false, res.Err()
	})
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}
