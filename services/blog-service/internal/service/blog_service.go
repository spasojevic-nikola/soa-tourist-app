package service

import (
	"context"
	"errors"
	"time"

	"blog-service/internal/dto"
	"blog-service/internal/models"
	"blog-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BlogService sadrži reference na repository.
type BlogService struct {
	Repo repository.BlogRepository
}

// NewBlogService kreira novu instancu BlogService-a.
func NewBlogService(repo repository.BlogRepository) *BlogService {
	return &BlogService{Repo: repo}
}

// CreateBlog kreira novi blog.
func (s *BlogService) CreateBlog(ctx context.Context, req dto.CreateBlogRequest, authorID uint) (*models.Blog, error) {
	blog := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  authorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Images:    req.Images,
		Comments:  []models.Comment{},
		Likes:     []uint{},
	}

	if err := s.Repo.CreateBlog(ctx, blog); err != nil {
		return nil, errors.New("failed to save blog to database")
	}
	return blog, nil
}

// AddComment dodaje komentar u blog.
func (s *BlogService) AddComment(ctx context.Context, blogID primitive.ObjectID, req dto.AddCommentRequest, authorID uint) (*models.Comment, error) {
	newComment := models.Comment{
		AuthorID:  authorID,
		Text:      req.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	update := bson.M{"$push": bson.M{"comments": newComment}}

	if err := s.Repo.UpdateBlog(ctx, blogID, update); err != nil {
		return nil, errors.New("failed to add comment to blog")
	}

	return &newComment, nil
}

// ToggleLike dodaje ili uklanja like korisnika.
func (s *BlogService) ToggleLike(ctx context.Context, blogID primitive.ObjectID, userID uint) (string, error) {
	blog, err := s.Repo.GetBlogByID(ctx, blogID)
	if err != nil || blog == nil {
		return "", errors.New("blog not found or internal error")
	}

	alreadyLiked := false
	for _, id := range blog.Likes {
		if id == userID {
			alreadyLiked = true
			break
		}
	}

	var update bson.M
	message := ""
	if alreadyLiked {
		update = bson.M{"$pull": bson.M{"likes": userID}}
		message = "Blog unliked successfully"
	} else {
		update = bson.M{"$addToSet": bson.M{"likes": userID}}
		message = "Blog liked successfully"
	}

	if err := s.Repo.UpdateBlog(ctx, blogID, update); err != nil {
		return "", errors.New("failed to update like status in database")
	}

	return message, nil
}

// GetAllBlogs vraća sve blogove.
func (s *BlogService) GetAllBlogs(ctx context.Context) ([]models.Blog, error) {
	return s.Repo.GetAll(ctx)
}

// GetBlogByID vraća blog po ID-ju.
func (s *BlogService) GetBlogByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error) {
	return s.Repo.GetByID(ctx, id)
}
