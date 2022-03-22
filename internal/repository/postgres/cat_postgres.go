package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/malkev1ch/first-task/internal/domain"
	"github.com/sirupsen/logrus"
	"strings"
)

type CatPostgres struct {
	db  *pgx.Conn
	ctx context.Context
}

func NewCatPostgres(db *pgx.Conn) *CatPostgres {
	return &CatPostgres{db: db}
}

func (r *CatPostgres) CreateCat(input domain.CreateCat) (*int, error) {
	insertCatQuery := fmt.Sprintf("INSERT INTO cats(name, date_birth, vaccinated) VALUES ($1, $2, $3) RETURNING id")
	var id int
	if err := r.db.QueryRow(r.ctx, insertCatQuery, input.Name, input.DateBirth, input.Vaccinated).Scan(&id); err != nil {
		logrus.Error(err, "Error occurred while inserting new row in table cats")
		return nil, err
	}

	return &id, nil
}

func (r *CatPostgres) GetCat(id int) (*domain.Cat, error) {
	getCatQuery := fmt.Sprintf("SELECT id, name, date_birth, vaccinated FROM cats WHERE id = $1")
	var cat domain.Cat
	if err := r.db.QueryRow(r.ctx, getCatQuery, id).Scan(&cat.Name, &cat.DateBirth, &cat.Vaccinated); err != nil {
		logrus.Error(err, "Error occurred while selecting row from table cats")
		return nil, err
	}

	return nil, nil
}

func (r *CatPostgres) UpdateCat(id int, input domain.UpdateCat) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.DateBirth != nil {
		setValues = append(setValues, fmt.Sprintf("date_birth=$%d", argId))
		args = append(args, *input.DateBirth)
		argId++
	}

	if input.Vaccinated != nil {
		setValues = append(setValues, fmt.Sprintf("vaccinated=$%d", argId))
		args = append(args, *input.Vaccinated)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	updateCatQuery := fmt.Sprintf("UPDATE cats SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)

	if _, err := r.db.Exec(r.ctx, updateCatQuery, args...); err != nil {
		logrus.Error(err, "Error occurred while updating row from table cats")
		return err
	}

	return nil
}

func (r *CatPostgres) DeleteCat(id int) error {
	deleteCatQuery := fmt.Sprintf("DELETE FROM cats WHERE id = $1")
	if _, err := r.db.Exec(r.ctx, deleteCatQuery, id); err != nil {
		logrus.Error(err, "Error occurred while deleting row from table cats")
		return err
	}

	return nil
}
