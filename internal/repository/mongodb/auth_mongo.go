package mongodb

import (
	"context"

	"github.com/malkev1ch/first-task/internal/repository"
)

func (r RepositoryMongo) CreateUser(ctx context.Context, input *repository.CreateUserInput) error {
	return nil
}

func (r RepositoryMongo) GetUserHashedPassword(ctx context.Context, email string) (string, string, error) {
	return "", "", nil
}

func (r RepositoryMongo) UpdateUserRefreshToken(ctx context.Context, refreshToken, id string) error {
	return nil
}

func (r RepositoryMongo) GetUserRefreshToken(ctx context.Context, id string) (string, error) {
	return "", nil
}
