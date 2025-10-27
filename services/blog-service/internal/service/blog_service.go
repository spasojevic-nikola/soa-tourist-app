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
	"go.mongodb.org/mongo-driver/mongo"
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

	var authorUsername string

	authURL := fmt.Sprintf("http://auth-service:8084/api/v1/auth/user/%d", authorID)
	resp, err := http.Get(authURL)
    
	if err != nil {
		log.Printf("Warning: Auth service unreachable while creating blog for ID %d: %v", authorID, err)
		authorUsername = "Unknown Author" // Postavi default
	} else {
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            var userData struct {
                Username string `json:"username"`
            }
            if err := json.NewDecoder(resp.Body).Decode(&userData); err == nil {
                authorUsername = userData.Username
            } else {
                log.Printf("Warning: Failed to parse username from auth response (%s) for ID %d: %v", resp.Status, authorID, err)
                authorUsername = "Unknown Author"
            }
        } else {
            log.Printf("Warning: Auth service returned status %s for ID %d", resp.Status, authorID)
            authorUsername = "Unknown Author"
        }
	}
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
		AuthorUsername: authorUsername,
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

	authURL := fmt.Sprintf("http://auth-service:8084/api/v1/auth/user/%d", authorID)
    resp, err := http.Get(authURL)
    if err != nil || resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch username from auth service")
    }
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("auth service returned status: %s", resp.Status)
	}

    var userData struct {
        Username string `json:"username"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
        return nil, fmt.Errorf("failed to parse username from auth response: %w", err)
    }

	newComment := models.Comment{
		ID:        primitive.NewObjectID(),
		AuthorID:  authorID,
		AuthorUsername: userData.Username,
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
    // 1. Ako je korisnik anoniman (ID = 0), vrati sve blogove
    if userID == 0 {
        log.Println("Anonymous user: returning all blogs")
        return s.Repo.GetAll(ctx)
    }

    // 2. KREIRANJE HTTP ZAHTEVA KA FOLLOWER SERVICE-u
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

    // 3. PROSLEĐIVANJE IDENTITETA KORISNIKA
    // Follower service ocekuje X-User-ID header koji postavlja API Gateway
    // Blog Service mora da prosledi ovaj identitet.
    req.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))

    // 4. SLANJE ZAHTEVA I OBRADA ODGOVORA
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        // Ako follower service nije dostupan, vrati samo blogove korisnika
        log.Printf("Follower service unavailable for user %d, returning only user's blogs: %v", userID, err)
        return s.Repo.GetBlogsByAuthorIDs(ctx, []uint{userID})
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Ako follower service vrati error, vrati samo blogove korisnika
        log.Printf("Follower service returned status %s for user %d, returning only user's blogs", resp.Status, userID)
        return s.Repo.GetBlogsByAuthorIDs(ctx, []uint{userID})
    }

    var followedIDs []uint
    if err := json.NewDecoder(resp.Body).Decode(&followedIDs); err != nil {
        // Ako ne možemo da dekodiramo odgovor, vrati samo blogove korisnika
        log.Printf("Failed to decode follower service response for user %d: %v, returning only user's blogs", userID, err)
        return s.Repo.GetBlogsByAuthorIDs(ctx, []uint{userID})
    }

    // 5. UKLJUČIMO I BLOGOVE SAMOG KORISNIKA
    // Korisnik uvek treba da vidi i svoje blogove na feed-u.
    followedIDs = append(followedIDs, userID)

    // 6. Ako korisnik ne prati nikoga, vrati samo njegove blogove
    // (ovo je redundantno sa linijom iznad ali ostavljamo za bezbednost)
    if len(followedIDs) == 1 && followedIDs[0] == userID {
        log.Printf("User %d doesn't follow anyone, returning only their blogs", userID)
    }

    // 7. POZIV REPOSITORY-JA SA LISTOM ID-JEVA
    return s.Repo.GetBlogsByAuthorIDs(ctx, followedIDs)
}

// UpdateComment ažurira tekst komentara (samo autor komentara) koristeći ID-je.
func (s *BlogService) UpdateComment(ctx context.Context, blogID primitive.ObjectID, commentID primitive.ObjectID, req dto.UpdateCommentRequest, userID uint) (*models.Comment, error) {
	// 1. DOHVATANJE BLOGA ZA PROVERU AUTORIZACIJE
	blog, err := s.Repo.GetBlogByID(ctx, blogID)
	if err != nil || blog == nil {
		// Provera da li je greška interna ili da blog ne postoji
		if err == mongo.ErrNoDocuments { // Pretpostavljajući da GetBlogByID vraća mongo.ErrNoDocuments
			return nil, errors.New("blog not found")
		}
		return nil, errors.New("failed to retrieve blog")
	}

	// 2. PRONALAŽENJE KOMENTARA I PROVERA VLASNIŠTVA
	var targetComment *models.Comment
	for i := range blog.Comments {
		if blog.Comments[i].ID == commentID {
			targetComment = &blog.Comments[i]
			break
		}
	}

	if targetComment == nil {
		return nil, errors.New("comment not found")
	}

	// PROVERA AUTORIZACIJE: Samo autor može menjati komentar
	if targetComment.AuthorID != userID {
		return nil, errors.New("unauthorized: only the comment author can update it")
	}

	// 3. KREIRANJE FILTERA I UPDATE DOKUMENTA ZA MONGODB
	updatedTime := time.Now()

	// Filter pronalazi dokument bloga po ID-ju I element u nizu komentara po njegovom ID-ju
	filter := bson.M{
		"_id":          blogID,
		"comments._id": commentID,
	}
	// Update koristi positional operator `$` da bi ažurirao element koji je pronađen filterom
	update := bson.M{
		"$set": bson.M{
			"comments.$.text":      req.Text,
			"comments.$.updatedAt": updatedTime, // Ažuriranje vremena izmene
		},
	}

	// 4. AŽURIRANJE U BAZI KORIŠĆENJEM GENERIČKE METODE REPOSITORY-JA
	if err := s.Repo.UpdateOne(ctx, filter, update); err != nil {
		return nil, fmt.Errorf("failed to update comment in database: %w", err)
	}

	// 5. AŽURIRANJE LOKALNOG OBJEKTA ZA POVRATAK KLIJENTU
	targetComment.Text = req.Text
	targetComment.UpdatedAt = updatedTime

	return targetComment, nil
}

// UpdateBlog ažurira blog post (samo autor).
func (s *BlogService) UpdateBlog(ctx context.Context, blogID primitive.ObjectID, req dto.UpdateBlogRequest, userID uint) (*models.Blog, error) {
	// 1. DOHVATANJE BLOGA ZA PROVERU AUTORIZACIJE
	blog, err := s.Repo.GetBlogByID(ctx, blogID)
	if err != nil || blog == nil {
		// Provera da li je greška interna ili da blog ne postoji
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("blog not found")
		}
		return nil, errors.New("failed to retrieve blog")
	}

	// PROVERA AUTORIZACIJE: Samo autor može menjati blog
	if blog.AuthorID != userID {
		return nil, errors.New("unauthorized: only the author can update the blog")
	}

	// 2. KONVERZIJA MARKDOWN-a U HTML (kao kod kreiranja)
	rawMarkdown := []byte(req.Content)
	
	// Konfiguracija HTML renderera (Standardne opcije + otvaranje linkova u novom tabu)
	opts := mdhtml.RendererOptions{Flags: mdhtml.CommonFlags | mdhtml.HrefTargetBlank}
	renderer := mdhtml.NewRenderer(opts)
	
	// Generisanje HTML-a
	htmlOutput := md.ToHTML(rawMarkdown, nil, renderer)
	
	currentTime := time.Now()

	// 3. KREIRANJE UPDATE DOKUMENTA
	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"content":     req.Content,
			"htmlContent": string(htmlOutput), // Čuvamo generisani HTML
			"images":      req.Images,
			"updatedAt":   currentTime,      // Ažuriranje vremena izmene
		},
	}

	// 4. AŽURIRANJE U BAZI KORIŠĆENJEM REPOSITORY-JA
	if err := s.Repo.UpdateBlog(ctx, blogID, update); err != nil {
		return nil, errors.New("failed to update blog in database")
	}
	
	// 5. VRAĆANJE AŽURIRANOG OBJEKTA
	
	// U idealnom slučaju, ažurirali bismo lokalni objekt, ali da bismo bili 100% sigurni
	// da je sve u bazi ispravno, najbolje je ponovo ga učitati.
	return s.Repo.GetBlogByID(ctx, blogID)
}