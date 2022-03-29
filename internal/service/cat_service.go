package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

func (s Service) Create(ctx context.Context, cat *model.Cat) (string, error) {
	id := uuid.New().String()
	cat.ID = id
	if err := s.redis.Cat.Set(ctx, cat); err != nil {
		return "", err
	}

	if err := s.repo.Create(ctx, cat); err != nil {
		return "", err
	}
	return id, nil
}

func (s Service) Get(ctx context.Context, id string) (*model.Cat, error) {
	cat, ex := s.redis.Cat.Get(ctx, id)
	if !ex {
		logrus.Info("got cat from database")
		return s.repo.Cat.Get(ctx, id)
	}

	logrus.Info("got cat from cache")
	return cat, nil
}

func (s Service) Update(ctx context.Context, id string, input *model.UpdateCat) (*model.Cat, error) {
	cat, err := s.repo.Cat.Update(ctx, id, input)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Cat.Set(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	if err := s.redis.Cat.Delete(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Cat.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s Service) UploadImage(ctx context.Context, id, path string) error {
	cat, err := s.repo.Cat.UploadImage(ctx, id, path)
	if err != nil {
		return err
	}

	if err := s.redis.Cat.Set(ctx, cat); err != nil {
		return err
	}

	return nil
}
