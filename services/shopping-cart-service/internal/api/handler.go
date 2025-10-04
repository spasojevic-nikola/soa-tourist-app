package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"shopping-cart-service/internal/dto"
	"shopping-cart-service/internal/service"
)

// Handler sadrži referencu na CartService
type Handler struct {
	Service *service.CartService
}

//  kreira novu instancu Handler-a
func NewHandler(service *service.CartService) *Handler {
	return &Handler{Service: service}
}

// za dobijanje trenutne korpe
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)

	cart, err := h.Service.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

//  za dodavanje ture u korpu
func (h *Handler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	var req dto.AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := GetUserID(r)

	cart, err := h.Service.AddItemToCart(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// Checkout endpoint za finalizaciju kupovine
func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)

	// U realnosti bi ovde išla validacija placanja??
	
	resp, err := h.Service.Checkout(r.Context(), userID)
	if err != nil {
		if err.Error() == "shopping cart is empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Checkout failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r)
    vars := mux.Vars(r)
    tourID := vars["tourId"] // ID ture iz URL-a

    if tourID == "" {
        http.Error(w, "Tour ID is required", http.StatusBadRequest)
        return
    }

    cart, err := h.Service.RemoveItem(r.Context(), userID, tourID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cart)
}

func (h *Handler) HasPurchaseToken(w http.ResponseWriter, r *http.Request) {
    userID := GetUserID(r)
    vars := mux.Vars(r)
    tourID := vars["tourId"]

    if tourID == "" {
        http.Error(w, "Tour ID is required", http.StatusBadRequest)
        return
    }

    hasPurchased, err := h.Service.HasPurchaseToken(r.Context(), userID, tourID)
    if err != nil {
        http.Error(w, "Failed to check purchase status", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"isPurchased": hasPurchased})
}