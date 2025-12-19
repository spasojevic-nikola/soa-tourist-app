package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"soa-tourist-app/follower-service/internal/service"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Service *service.FollowerService
}

func NewHandler(s *service.FollowerService) *Handler {
	return &Handler{Service: s}
}

// Follow je HTTP handler za pracenje korisnika
func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"endpoint": "/api/followers/follow/{id}",
		"method":   "POST",
		"ip":       r.RemoteAddr,
	}).Info("Follow request received")

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
		log.WithFields(log.Fields{
			"endpoint":    "/api/followers/follow/{id}",
			"followed_id": followedIdStr,
			"error":       err.Error(),
		}).Error("Invalid user ID format")

		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Pozovi servis da odradi logiku
	err = h.Service.Follow(followerId, uint(followedId))
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":    "/api/followers/follow/{id}",
			"follower_id": followerId,
			"followed_id": followedId,
			"error":       err.Error(),
		}).Error("Failed to create follow relationship")
		// Vrati gresku ako npr. korisnik pokusa da zaprati sam sebe
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithFields(log.Fields{
		"endpoint":    "/api/followers/follow/{id}",
		"follower_id": followerId,
		"followed_id": followedId,
	}).Info("User followed successfully")

	// Vrati uspešan odgovor
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CheckFollow(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"endpoint": "/api/followers/check-follow/{id}",
		"method":   "GET",
		"ip":       r.RemoteAddr,
	}).Info("Check follow request received")

	followerId, _ := r.Context().Value(userKey).(uint)
	vars := mux.Vars(r)
	followedId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/followers/check-follow/{id}",
			"error":    err.Error(),
		}).Error("Invalid user ID in check follow")

		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	follows, err := h.Service.CheckFollows(followerId, uint(followedId))
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":    "/api/followers/check-follow/{id}",
			"follower_id": followerId,
			"followed_id": followedId,
			"error":       err.Error(),
		}).Error("Failed to check follow status")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.WithFields(log.Fields{
		"endpoint":    "/api/followers/check-follow/{id}",
		"follower_id": followerId,
		"followed_id": followedId,
		"follows":     follows,
	}).Info("Follow status checked successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"follows": follows})
}

func (h *Handler) Unfollow(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{
		"endpoint": "/api/followers/unfollow/{id}",
		"method":   "DELETE",
		"ip":       r.RemoteAddr,
	}).Info("Unfollow request received")

	followerId, ok := r.Context().Value(userKey).(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	followedId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/followers/unfollow/{id}",
			"error":    err.Error(),
		}).Error("Invalid user ID in unfollow")

		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.Service.Unfollow(followerId, uint(followedId))
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint":    "/api/followers/unfollow/{id}",
			"follower_id": followerId,
			"followed_id": followedId,
			"error":       err.Error(),
		}).Error("Failed to remove follow relationship")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"endpoint":    "/api/followers/unfollow/{id}",
		"follower_id": followerId,
		"followed_id": followedId,
	}).Info("User unfollowed successfully")

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Izvuci ID ulogovanog korisnika iz konteksta
	log.WithFields(log.Fields{
		"endpoint": "/api/followers/recommendations",
		"method":   "GET",
		"ip":       r.RemoteAddr,
	}).Info("Get recommendations request received")

	currentUserID, ok := r.Context().Value(userKey).(uint)
	if !ok || currentUserID == 0 { // Proverite da li je ID validan
		http.Error(w, "Unauthorized or missing user ID", http.StatusUnauthorized)
		return
	}

	recommendations, err := h.Service.GetRecommendations(currentUserID)
	if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/followers/recommendations",
			"user_id":  currentUserID,
			"error":    err.Error(),
		}).Error("Failed to fetch recommendations")

		http.Error(w, "Error fetching recommendations: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.WithFields(log.Fields{
		"endpoint": "/api/followers/recommendations",
		"user_id":  currentUserID,
		"count":    len(recommendations),
	}).Info("Recommendations fetched successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}

func (h *Handler) GetFollowingIDs(w http.ResponseWriter, r *http.Request) {
    // Izvuci ID ulogovanog korisnika iz konteksta
	log.WithFields(log.Fields{
		"endpoint": "/api/followers/following",
		"method":   "GET",
		"ip":       r.RemoteAddr,
	}).Info("Get following list request received")

    currentUserID, ok := r.Context().Value(userKey).(uint)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    ids, err := h.Service.GetFollowingIDs(currentUserID)
    if err != nil {
		log.WithFields(log.Fields{
			"endpoint": "/api/followers/following",
			"user_id":  currentUserID,
			"error":    err.Error(),
		}).Error("Failed to get following list")

        http.Error(w, "Failed to get following list", http.StatusInternalServerError)
        return
    }

	log.WithFields(log.Fields{
		"endpoint": "/api/followers/following",
		"user_id":  currentUserID,
		"count":    len(ids),
	}).Info("Following list retrieved successfully")
	
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ids)
}