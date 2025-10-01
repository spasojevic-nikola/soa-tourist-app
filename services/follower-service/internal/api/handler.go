package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"soa-tourist-app/follower-service/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	Service *service.FollowerService
}

func NewHandler(s *service.FollowerService) *Handler {
	return &Handler{Service: s}
}

// Follow je HTTP handler za pracenje korisnika
func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	// Izvuci ID ulogovanog korisnika iz konteksta (ovo će postaviti middleware)
	followerId, ok := r.Context().Value(userKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Izvuci ID korisnika koga treba zapratiti iz URL-a
	vars := mux.Vars(r)
	followedIdStr := vars["id"]
	followedId, err := strconv.ParseUint(followedIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Pozovi servis da odradi logiku
	err = h.Service.Follow(followerId, uint(followedId))
	if err != nil {
		// Vrati gresku ako npr. korisnik pokusa da zaprati sam sebe
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Vrati uspešan odgovor
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CheckFollow(w http.ResponseWriter, r *http.Request) {
	followerId, _ := r.Context().Value(userKey).(uint)
	vars := mux.Vars(r)
	followedId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	follows, err := h.Service.CheckFollows(followerId, uint(followedId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"follows": follows})
}

func (h *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {
	followerId, ok := r.Context().Value(userKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.Service.Unfollow(followerId, uint(followedId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Izvuci ID ulogovanog korisnika iz konteksta
	currentUserID, ok := r.Context().Value(userKey).(uint)
	if !ok || currentUserID == 0 { // Proverite da li je ID validan
		http.Error(w, "Unauthorized or missing user ID", http.StatusUnauthorized)
		return
	}

	recommendations, err := h.Service.GetRecommendations(currentUserID)
	if err != nil {
		http.Error(w, "Error fetching recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}