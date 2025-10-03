package api

import (
	"encoding/json"
	"net/http"

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