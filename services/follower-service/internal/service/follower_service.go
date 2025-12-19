package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"errors"
	"net/http"
	"time"    
	"soa-tourist-app/follower-service/internal/repository"
	"soa-tourist-app/follower-service/internal/dto"


)

type FollowerService struct {
	Repo *repository.FollowerRepository
	HttpClient *http.Client 
}
// Pomoćna struktura za dekodiranje odgovora iz stakeholders servisa
type StakeholderUser struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	ProfileImage string `json:"profile_image"`
}

// NewFollowerService kreira novu instancu servisa
func NewFollowerService(repo *repository.FollowerRepository) *FollowerService {
	return &FollowerService{Repo: repo,
						   HttpClient: &http.Client{Timeout: 10 * time.Second}, 
	}
}

// Follow proverava logiku i poziva repozitorijum da zaprati korisnika
func (s *FollowerService) Follow(followerId, followedId uint) error {
	// Primer biznis logike: ne dozvoli korisniku da zaprati sam sebe
	if followerId == followedId {
		return errors.New("cannot follow yourself")
	}

	return s.Repo.Follow(followerId, followedId)
}

func (s *FollowerService) Unfollow(followerId, followedId uint) error {
	return s.Repo.Unfollow(followerId, followedId)
}

func (s *FollowerService) CheckFollows(followerId, followedId uint) (bool, error) {
	return s.Repo.CheckFollows(followerId, followedId)
}

func (s *FollowerService) GetRecommendations(currentUserID uint) ([]dto.RecommendationDTO, error) {
	recommendedUsers, err := s.Repo.GetRecommendations(currentUserID)
	if err != nil {
		return nil, err
	}
	if len(recommendedUsers) == 0 {
		return []dto.RecommendationDTO{}, nil
	}

	var ids []string
	for _, rec := range recommendedUsers {
		ids = append(ids, fmt.Sprint(rec.UserID))
	}
	idsQueryParam := strings.Join(ids, ",")

	stakeholdersURL := fmt.Sprintf("http://stakeholders-service:8080/api/v1/users/batch?ids=%s", idsQueryParam)
	resp, err := s.HttpClient.Get(stakeholdersURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userProfiles []StakeholderUser
	if err := json.NewDecoder(resp.Body).Decode(&userProfiles); err != nil {
		return nil, err
	}

	profilesMap := make(map[uint]StakeholderUser)
	for _, profile := range userProfiles {
		profilesMap[profile.ID] = profile
	}

	var finalRecommendations []dto.RecommendationDTO
	for _, rec := range recommendedUsers {
		if profile, ok := profilesMap[rec.UserID]; ok {
			finalRecommendations = append(finalRecommendations, dto.RecommendationDTO{
				UserID:       profile.ID,
				Username:     profile.Username,
				FirstName:    profile.FirstName,
				LastName:     profile.LastName,
				ProfileImage: profile.ProfileImage,
				Score:        rec.Score,
			})
		}
	}

	return finalRecommendations, nil
}
func (s *FollowerService) GetFollowingIDs(followerId uint) ([]uint, error) {
    return s.Repo.GetFollowingIDs(followerId)
}