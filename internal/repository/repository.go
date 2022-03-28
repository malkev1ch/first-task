package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/malkev1ch/first-task/internal/model"
)

type CreateUserInput struct {
	ID           string
	UserName     string
	Email        string
	Password     string
	RefreshToken string
}

type Cat interface {
	Create(ctx context.Context, cat *model.Cat) error
	Get(ctx context.Context, id string) (*model.Cat, error)
	Update(ctx context.Context, id string, input *model.UpdateCat) error
	Delete(ctx context.Context, id string) error
	UploadImage(ctx context.Context, id string, path string) error
}

type Auth interface {
	CreateUser(ctx context.Context, input *CreateUserInput) error
	GetUserHashedPassword(ctx context.Context, email string) (string, string, error)
	UpdateUserRefreshToken(ctx context.Context, refreshToken, id string) error
	GetUserRefreshToken(ctx context.Context, id string) (string, error)
}

type Repository struct {
	Cat
	Auth
}

func NewRepositoryPostgres(DB *pgxpool.Pool) *Repository {
	return &Repository{
		Cat:  NewCatRepository(DB),
		Auth: NewAuthRepository(DB),
	}
}

func NewRepositoryMongo(DB *mongo.Client) *Repository {
	return &Repository{
		Cat:  NewCatRepositoryMongo(DB),
		Auth: NewAuthRepositoryMongo(DB),
	}
}
