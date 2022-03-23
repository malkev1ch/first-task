package repository

import (
	"github.com/jackc/pgx/v4"
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/malkev1ch/first-task/internal/repository/postgres"
)

type Cat interface {
	CreateCat(input domain.CreateCat) (*int, error)
	GetCat(id int) (*domain.Cat, error)
	UpdateCat(id int, input domain.UpdateCat) error
	DeleteCat(id int) error
	UploadImage(id int, path string) error
}

type Repository struct {
	Cat
}

func NewRepositoryPostgres(db *pgx.Conn) *Repository {
	return &Repository{
		Cat: postgres.NewCatPostgres(db),
	}
}
