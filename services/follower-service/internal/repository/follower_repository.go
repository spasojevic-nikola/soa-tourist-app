package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"soa-tourist-app/follower-service/internal/models"
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

func (repo *FollowerRepository) GetRecommendations(currentUserID uint) ([]models.RecommendationModel, error) {
	ctx := context.Background()
	session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
        MATCH (ja:User {id: $currentUserID})-[:FOLLOWS]->(prijatelj:User)
        MATCH (prijatelj)-[:FOLLOWS]->(kandidat:User)
        WHERE NOT (ja)-[:FOLLOWS]->(kandidat)
          AND ja.id <> kandidat.id
        RETURN kandidat.id AS recommendedUserID, 
               count(prijatelj) AS score
        ORDER BY score DESC
        LIMIT 10
    `
	params := map[string]interface{}{
		"currentUserID": currentUserID,
	}

	recommendations := make([]models.RecommendationModel, 0)

	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			
			userID, ok := record.Get("recommendedUserID")
			if !ok {
				continue
			}
			score, _ := record.Get("score")

			recommendations = append(recommendations, models.RecommendationModel{
				UserID: uint(userID.(int64)), 
				Score:  int(score.(int64)),
			})
		}
		if result.Err() != nil {
			return nil, result.Err()
		}
		return nil, nil
	})

	return recommendations, err
}

//  vraca listu ID-jeva korisnika koje dati korisnik prati
func (repo *FollowerRepository) GetFollowingIDs(followerId uint) ([]uint, error) {
    ctx := context.Background()
    session := repo.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
    defer session.Close(ctx)

    query := `
        MATCH (follower:User {id: $followerId})-[:FOLLOWS]->(followed:User)
        RETURN followed.id AS followedId
    `
    params := map[string]interface{}{
        "followerId": followerId,
    }

    followingIDs := make([]uint, 0)

    result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
        res, err := tx.Run(ctx, query, params)
        if err != nil {
            return nil, err
        }

        for res.Next(ctx) {
            record := res.Record()
            id, ok := record.Get("followedId")
            if ok {
                followingIDs = append(followingIDs, uint(id.(int64)))
            }
        }
        return nil, res.Err()
    })

    if err != nil {
        return nil, err
    }
    
    // Provera da li je result nil, iako ExecuteRead vraca (interface{}, error)
    if result != nil {
    }

    return followingIDs, nil
}