package service

import (
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/malkev1ch/first-task/internal/repository"
)

type Cat interface {
	CreateCat(input domain.CreateCat) (*int, error)
	GetCat(id int) (*domain.Cat, error)
	UpdateCat(id int, input domain.UpdateCat) error
	DeleteCat(id int) error
	UploadImage(id int, path string) error
}

type Service struct {
	Cat
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Cat: NewCatService(repo.Cat),
	}
}
