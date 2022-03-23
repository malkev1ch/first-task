package repository

import (
	"context"
	"github.com/malkev1ch/first-task/internal/model"
)

type Repository interface {
	Create(ctx context.Context, cat *model.Cat) (int, error)
	Get(ctx context.Context, id int) (*model.Cat, error)
	Update(ctx context.Context, id int, input *model.Cat) error
	Delete(ctx context.Context, id int) error
	UploadImage(ctx context.Context, id int, path string) error
}
