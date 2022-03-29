package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepositoryMongo type represents postgres behavior for authentication.
type AuthRepositoryMongo struct {
	DB *mongo.Client
}

func NewAuthRepositoryMongo(db *mongo.Client) *AuthRepositoryMongo {
	return &AuthRepositoryMongo{
		DB: db,
	}
}

func (r AuthRepositoryMongo) CreateUser(ctx context.Context, input *CreateUserInput) error {
	return nil
}

func (r AuthRepositoryMongo) GetUserHashedPassword(ctx context.Context, email string) (string, string, error) {
	return "", "", nil
}

func (r AuthRepositoryMongo) UpdateUserRefreshToken(ctx context.Context, refreshToken, id string) error {
	return nil
}

func (r AuthRepositoryMongo) GetUserRefreshToken(ctx context.Context, id string) (string, error) {
	return "", nil
}
