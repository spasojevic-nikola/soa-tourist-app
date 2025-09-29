package service

import (
	"errors"
	"soa-tourist-app/follower-service/internal/repository"
)

type FollowerService struct {
	Repo *repository.FollowerRepository
}

// NewFollowerService kreira novu instancu servisa
func NewFollowerService(repo *repository.FollowerRepository) *FollowerService {
	return &FollowerService{Repo: repo}
}

// Follow proverava logiku i poziva repozitorijum da zaprati korisnika
func (s *FollowerService) Follow(followerId, followedId uint) error {
	// Primer biznis logike: ne dozvoli korisniku da zaprati sam sebe
	if followerId == followedId {
		return errors.New("cannot follow yourself")
	}

	return s.Repo.Follow(followerId, followedId)
}