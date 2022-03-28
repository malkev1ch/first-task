package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

func (s Service) Create(ctx context.Context, cat *model.Cat) (string, error) {
	id := uuid.New().String()
	cat.ID = id

	if err := s.redis.Cat.Save(ctx, cat); err != nil {
		logrus.Error(err, "service: error occurred while setting data in redis")
		return "", fmt.Errorf("service: error occurred while setting data in redis - %w", err)
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return "", err
	}
	return id, nil
}

func (s Service) Get(ctx context.Context, id string) (*model.Cat, error) {
	return s.repo.Cat.Get(ctx, id)
}

func (s Service) Update(ctx context.Context, id string, input *model.UpdateCat) error {
	return s.repo.Cat.Update(ctx, id, input)
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.repo.Cat.Delete(ctx, id)
}

func (s Service) UploadImage(ctx context.Context, id, path string) error {
	return s.repo.Cat.UploadImage(ctx, id, path)
}
