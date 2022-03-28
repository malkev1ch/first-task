package service

import (
	"github.com/malkev1ch/first-task/internal/rediscache"
	"github.com/malkev1ch/first-task/internal/repository"
)

type Service struct {
	repo  *repository.Repository
	redis *rediscache.Cache
}

func NewService(repo *repository.Repository, redis *rediscache.Cache) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
	}
}
