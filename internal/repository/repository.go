package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/malkev1ch/first-task/internal/repository/postgres"
	"time"
)

type Cat interface {
	CreateCat(input domain.CreateCat) (*int, error)
	GetCat(id int) (*domain.Cat, error)
	UpdateCat(id int, input domain.UpdateCat) error
	DeleteCat(id int) error
}

type Repository struct {
	Cat
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRepositoryPostgres(db *pgx.Conn) *Repository {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	return &Repository{
		Cat:    postgres.NewCatPostgres(db),
		ctx:    ctx,
		cancel: cancel,
	}
}
