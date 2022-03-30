package service

import (
	"context"

	"github.com/malkev1ch/first-task/internal/model"
	"github.com/malkev1ch/first-task/internal/rediscache"
	"github.com/malkev1ch/first-task/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/service_mock.go

type Cat interface {
	Create(ctx context.Context, cat *model.Cat) (string, error)
	Get(ctx context.Context, id string) (*model.Cat, error)
	Update(ctx context.Context, id string, input *model.UpdateCat) (*model.Cat, error)
	Delete(ctx context.Context, id string) error
	UploadImage(ctx context.Context, id, path string) error
}

type Auth interface {
	SignUp(ctx context.Context, input *model.CreateUser) (*model.Tokens, error)
	SignIn(ctx context.Context, input *model.AuthUser) (*model.Tokens, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (*model.Tokens, error)
}

type Service struct {
	Cat
	Auth
}

func NewService(repo *repository.Repository, redis *rediscache.Cache) *Service {
	return &Service{
		Cat:  NewCatService(repo, redis),
		Auth: NewAuthService(repo),
	}
}
