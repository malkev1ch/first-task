package repository

import (
	"context"

	"github.com/malkev1ch/first-task/internal/model"
)

type Cat interface {
	Create(ctx context.Context, cat *model.CreateCat) (string, error)
	Get(ctx context.Context, id string) (model.Cat, error)
	Update(ctx context.Context, id string, input *model.UpdateCat) error
	Delete(ctx context.Context, id string) error
	UploadImage(ctx context.Context, id string, path string) error
}

type Repository interface {
	Cat
}
