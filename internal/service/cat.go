package service

import (
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/malkev1ch/first-task/internal/repository"
)

type CatService struct {
	repo repository.Cat
}

func NewCatService(repo repository.Cat) *CatService {
	return &CatService{repo: repo}
}

func (s *CatService) CreateCat(input domain.CreateCat) (*int, error) {
	return s.repo.CreateCat(input)
}

func (s *CatService) GetCat(id int) (*domain.Cat, error) {
	return s.repo.GetCat(id)
}

func (s *CatService) UpdateCat(id int, input domain.UpdateCat) error {
	return s.repo.UpdateCat(id, input)
}

func (s *CatService) DeleteCat(id int) error {
	return s.repo.DeleteCat(id)
}

func (s *CatService) UploadImage(id int, path string) error {
	return s.repo.UploadImage(id, path)
}

