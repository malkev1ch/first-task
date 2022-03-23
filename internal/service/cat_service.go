package service

import (
	"context"
	"github.com/malkev1ch/first-task/internal/model"
)

func (s Service) Create(ctx context.Context, cat *model.Cat) (int, error) {
	return s.repo.Create(ctx, cat)
}

func (s Service) Get(ctx context.Context, id int) (*model.Cat, error) {
	return s.repo.Get(ctx, id)
}

func (s Service) Update(ctx context.Context, id int, input *model.Cat) error {
	return s.repo.Update(ctx, id, input)
}

func (s Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s Service) UploadImage(ctx context.Context, id int, path string) error {
	return s.repo.UploadImage(ctx, id, path)
}
