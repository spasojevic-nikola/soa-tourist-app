package repository

import (
	"context"
	"blog-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// BlogRepository je interfejs koji definiše metode za rad sa blogovima.
type BlogRepository interface {
	CreateBlog(ctx context.Context, blog *models.Blog) error
	GetBlogByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error)
	UpdateBlog(ctx context.Context, id primitive.ObjectID, update bson.M) error
	GetAll(ctx context.Context) ([]models.Blog, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error)
	GetBlogsByAuthorIDs(ctx context.Context, authorIDs []uint) ([]models.Blog, error)
}

// mongoBlogRepository je konkretna implementacija BlogRepository koristeći MongoDB.
type mongoBlogRepository struct {
	collection *mongo.Collection
}

// NewBlogRepository kreira novi MongoDB blog repository.
func NewBlogRepository(db *mongo.Database) BlogRepository {
	return &mongoBlogRepository{
		collection: db.Collection("blogs"),
	}
}

// CreateBlog dodaje novi blog u bazu.
func (r *mongoBlogRepository) CreateBlog(ctx context.Context, blog *models.Blog) error {
	_, err := r.collection.InsertOne(ctx, blog)
	return err
}

// GetBlogByID vraća blog po ID-ju.
func (r *mongoBlogRepository) GetBlogByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error) {
	var blog models.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &blog, err
}

// UpdateBlog ažurira blog po ID-ju.
func (r *mongoBlogRepository) UpdateBlog(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// GetAll vraća sve blogove iz baze.
func (r *mongoBlogRepository) GetAll(ctx context.Context) ([]models.Blog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var blogs []models.Blog
	if err = cursor.All(ctx, &blogs); err != nil {
		return nil, err
	}
	return blogs, nil
}
//samo blogove korisnika koje pratim
func (r *mongoBlogRepository) GetBlogsByAuthorIDs(ctx context.Context, authorIDs []uint) ([]models.Blog, error) {
    // Kreiramo filter koji traži blogove gde je 'authorId' u nizu 'authorIDs'
    // Ovo je ekvivalent SQL-ovog "WHERE authorId IN (id1, id2, ...)"
    filter := bson.M{"authorId": bson.M{"$in": authorIDs}}

    cursor, err := r.collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var blogs []models.Blog
    if err = cursor.All(ctx, &blogs); err != nil {
        return nil, err
    }
    return blogs, nil
}

// GetByID vraća blog po ID-ju (slično GetBlogByID).
func (r *mongoBlogRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Blog, error) {
	var blog models.Blog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&blog)
	if err != nil {
		return nil, err
	}
	return &blog, nil
}
