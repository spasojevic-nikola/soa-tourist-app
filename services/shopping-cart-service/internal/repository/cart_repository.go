package repository

import (
	"context"
	"errors"
	"time"

	"shopping-cart-service/internal/models"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// interfejs za rad sa korpom i tokenima
type CartRepository interface {
	GetCartByUserID(ctx context.Context, userID uint) (*models.ShoppingCart, error)
	CreateCart(ctx context.Context, cart *models.ShoppingCart) error
	UpdateCart(ctx context.Context, cart *models.ShoppingCart) error
    RemoveItem(ctx context.Context, userID uint, tourID string) error 
	DeleteCart(ctx context.Context, userID uint) error 
	CreatePurchaseTokens(ctx context.Context, tokens []models.TourPurchaseToken) ([]primitive.ObjectID, error)
}

type mongoCartRepository struct {
	cartCollection    *mongo.Collection
	tokenCollection *mongo.Collection
}

// kreira novi MongoDB repository.
func NewCartRepository(db *mongo.Database) CartRepository {
	return &mongoCartRepository{
		cartCollection: db.Collection("shopping_carts"),
		tokenCollection: db.Collection("purchase_tokens"),
	}
}

// pronalazi korpu po ID-ju korisnika
func (r *mongoCartRepository) GetCartByUserID(ctx context.Context, userID uint) (*models.ShoppingCart, error) {
	var cart models.ShoppingCart
	err := r.cartCollection.FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Korpa ne postoji
	}
	return &cart, err
}

// kreira novu korpu
func (r *mongoCartRepository) CreateCart(ctx context.Context, cart *models.ShoppingCart) error {
	cart.Updated = time.Now()
	_, err := r.cartCollection.InsertOne(ctx, cart)
	return err
}

// azurira postojeću korpu (koristi ReplaceOne da bi zamenio ceo dokument novim stanjem)
func (r *mongoCartRepository) UpdateCart(ctx context.Context, cart *models.ShoppingCart) error {
	cart.Updated = time.Now()
	
	result, err := r.cartCollection.ReplaceOne(
		ctx,
		bson.M{"userId": cart.UserID},
		cart,
	)

	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return errors.New("cart not found or no changes made")
	}
	return nil
}

func (r *mongoCartRepository) DeleteCart(ctx context.Context, userID uint) error {
	_, err := r.cartCollection.DeleteOne(ctx, bson.M{"userId": userID})
	return err
}

func (r *mongoCartRepository) RemoveItem(ctx context.Context, userID uint, tourID string) error {
    update := bson.M{
        "$pull": bson.M{
            "items": bson.M{"tourId": tourID},
        },
    }
    
    // Tražimo korpu po UserID-ju
    _, err := r.cartCollection.UpdateOne(ctx, bson.M{"userId": userID}, update)
    return err
}

//  snima tokene kupljenih tura
func (r *mongoCartRepository) CreatePurchaseTokens(ctx context.Context, tokens []models.TourPurchaseToken) ([]primitive.ObjectID, error) {
	if len(tokens) == 0 {
		return nil, nil
	}

	docs := make([]interface{}, len(tokens))
	for i, t := range tokens {
		t.PurchaseTime = time.Now()
		docs[i] = t
	}

	result, err := r.tokenCollection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	
	objectIDs := make([]primitive.ObjectID, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		objectIDs[i] = id.(primitive.ObjectID)
	}
	return objectIDs, nil
}