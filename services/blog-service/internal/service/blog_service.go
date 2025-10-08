package service

import (
	"context"
	"errors"
	"time"
	"encoding/json"
	"fmt"     
	"net/http" 
	"log"
	"os"


	"blog-service/internal/dto"
	"blog-service/internal/models"
	"blog-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	md "github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	mdparser "github.com/gomarkdown/markdown/parser"
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

	 // 1. Definišemo ekstenzije koje želimo da naš parser podržava
	 extensions := mdparser.CommonExtensions | mdparser.AutoHeadingIDs | mdparser.Strikethrough
    
	 // 2. Kreiramo novi parser sa tim ekstenzijama
	 p := mdparser.NewWithExtensions(extensions)
	// 1. KONVERZIJA MARKDOWN-a U HTML
    rawMarkdown := []byte(req.Content)
    
    // Konfiguracija HTML renderera (Standardne opcije + otvaranje linkova u novom tabu)
    opts := mdhtml.RendererOptions{Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank}
    renderer := mdhtml.NewRenderer(opts)
    
    // Generisanje HTML-a
    htmlOutput := md.ToHTML(rawMarkdown, p, renderer)

	blog := &models.Blog{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		HTMLContent: string(htmlOutput), // Čuvamo generisani HTML
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

func (s *BlogService) GetFeedForUser(ctx context.Context, userID uint) ([]models.Blog, error) {
    // 1. KREIRANJE HTTP ZAHTEVA KA FOLLOWER SERVICE-u
    // U realnoj aplikaciji, URL bi bio u konfiguraciji (npr. env varijabla)
    //followerServiceURL := "http://follower-service:8080/api/followers/following"
	followerServiceBaseURL := os.Getenv("FOLLOWER_SERVICE_URL")

	if followerServiceBaseURL == "" {
        log.Fatal("FATAL: FOLLOWER_SERVICE_URL environment variable is not set.")
    }

    // Sastavljamo pun URL
    followerServiceURL := fmt.Sprintf("%s/api/followers/following", followerServiceBaseURL)

    req, err := http.NewRequestWithContext(ctx, "GET", followerServiceURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request to follower service: %w", err)
    }

    // 2. PROSLEĐIVANJE IDENTITETA KORISNIKA
    // Follower service ocekuje X-User-ID header koji postavlja API Gateway
    // Blog Service mora da prosledi ovaj identitet.
    req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

    // 3. SLANJE ZAHTEVA I OBRADA ODGOVORA
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        // Ovo se desava ako je Follower Service pao ili mreža ne radi
        return nil, fmt.Errorf("follower service is unavailable: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("follower service returned status: %s", resp.Status)
    }

    var followedIDs []uint
    if err := json.NewDecoder(resp.Body).Decode(&followedIDs); err != nil {
        return nil, fmt.Errorf("failed to decode response from follower service: %w", err)
    }

    // 4. UKLJUCjEM I BLOGOVE SAMOG KORISNIKA
    // Korisnik uvek treba da vidi i svoje blogove na feed-u.
    followedIDs = append(followedIDs, userID)

    // Ako korisnik ne prati nikoga, vrati samo njegove blogove
    if len(followedIDs) == 0 {
        followedIDs = []uint{userID}
    }

    // 5. POZIV REPOSITORY-JA SA LISTOM ID-JEVA
    return s.Repo.GetBlogsByAuthorIDs(ctx, followedIDs)
}
