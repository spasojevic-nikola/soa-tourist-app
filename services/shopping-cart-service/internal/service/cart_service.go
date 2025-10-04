package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"shopping-cart-service/internal/dto"
	"shopping-cart-service/internal/models"
	"shopping-cart-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CartService sadrži reference na repository.
type CartService struct {
	Repo repository.CartRepository
}

// NewCartService kreira novu instancu CartService-a.
func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{Repo: repo}
}

// calculateTotal računa ukupnu cenu svih stavki.
func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price
	}
	return total
}

// GetCart vraća korpu za datog korisnika (kreira je ako ne postoji).
func (s *CartService) GetCart(ctx context.Context, userID uint) (*models.ShoppingCart, error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping cart")
	}
	if cart == nil {
		// Ako korpa ne postoji, kreiraj novu praznu korpu
		cart = &models.ShoppingCart{
			UserID: userID,
			Items:  []models.OrderItem{},
			Total:  0,
			ID:     primitive.NewObjectID(),
		}
		if err := s.Repo.CreateCart(ctx, cart); err != nil {
			return nil, errors.New("failed to create new shopping cart")
		}
	}
	return cart, nil
}

// AddItemToCart dodaje stavku u korpu i ažurira total.
func (s *CartService) AddItemToCart(ctx context.Context, userID uint, req dto.AddItemRequest) (*models.ShoppingCart, error) {
	// Proveri dostupnost ture u Tour servisu ovde (nije implementirano)
	
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	newItem := models.OrderItem{
		TourID: req.TourID,
		Name:   req.Name,
		Price:  req.Price,
	}

	cart.Items = append(cart.Items, newItem)
	cart.Total = calculateTotal(cart.Items)

	if err := s.Repo.UpdateCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to update cart: %w", err)
	}

	return cart, nil
}

// Checkout obrađuje kupovinu: kreira tokene i briše korpu.
func (s *CartService) Checkout(ctx context.Context, userID uint) (*dto.TourPurchaseResponse, error) {
	cart, err := s.Repo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to retrieve shopping cart for checkout")
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("shopping cart is empty")
	}

	// 1. Kreiranje tokena za svaku stavku
	tokens := make([]models.TourPurchaseToken, len(cart.Items))
	for i, item := range cart.Items {
		tokens[i] = models.TourPurchaseToken{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			TourID: item.TourID,
		}
	}

	// 2. Snimanje tokena u bazu
	tokenIDs, err := s.Repo.CreatePurchaseTokens(ctx, tokens)
	if err != nil {
		return nil, errors.New("failed to create purchase tokens")
	}

	// 3. Brisanje korpe
	// Turista dobija tokene i korpa se briše.
	if err := s.Repo.DeleteCart(ctx, userID); err != nil {
		// Logovati, ali ne zaustavljati proces jer su tokeni snimljeni
		fmt.Printf("Warning: Failed to delete cart after successful purchase for user %d: %v\n", userID, err)
	}

	return &dto.TourPurchaseResponse{
		Tokens: tokenIDs,
		Message: fmt.Sprintf("Purchase successful. %d items bought for %.2f.", len(cart.Items), cart.Total),
	}, nil
}
func (s *CartService) RemoveItem(ctx context.Context, userID uint, tourID string) (*models.ShoppingCart, error) {
	// 1. Ukloni stavku iz baze
	if err := s.Repo.RemoveItem(ctx, userID, tourID); err != nil {
		log.Printf("ERROR: Failed to remove item %s from cart of User %d. Error: %v", tourID, userID, err)
		return nil, fmt.Errorf("failed to remove item from cart: %w", err)
	}

	// 2. Dohvati azuriranu korpu
	cart, err := s.Repo.GetCartByUserID(ctx, userID)
	if err != nil || cart == nil {
		return nil, errors.New("cart not found after removal, internal inconsistency")
	}

	// 3. Ponovo izracunaj total
	cart.Total = calculateTotal(cart.Items)
	cart.Updated = time.Now() // Koristimo time.Now()

	// 4. Ažuriraj total u bazi (koristeći ReplaceOne u UpdateCart)
	if err := s.Repo.UpdateCart(ctx, cart); err != nil {
		return nil, fmt.Errorf("failed to recalculate cart total: %w", err)
	}
    
	log.Printf("INFO: User %d removed item %s. New Total: %.2f", userID, tourID, cart.Total)
	return cart, nil
}
