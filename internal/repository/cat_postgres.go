package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

// CatRepository type represents postgres object cat structure and behavior.
type CatRepository struct {
	DB *pgxpool.Pool
}

func NewCatRepository(db *pgxpool.Pool) *CatRepository {
	return &CatRepository{
		DB: db,
	}
}

// Create method saves object Cat into postgres database.
func (r CatRepository) Create(ctx context.Context, input *model.Cat) error {
	logrus.WithFields(logrus.Fields{
		"id":         input.ID,
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Info("postgres repository: create cat")

	insertCatQuery := "INSERT INTO cats(id, name, date_birth, vaccinated) VALUES ($1, $2, $3, $4)"

	if _, err := r.DB.Exec(ctx, insertCatQuery, input.ID, input.Name, input.DateBirth, input.Vaccinated); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while inserting new row in table cats")
		return fmt.Errorf("postgres repository: can't create cat - %w", err)
	}
	return nil
}

// Get method returns object Cat from postgres database
// with selection by id.
func (r CatRepository) Get(ctx context.Context, id string) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Info("postgres repository: get cat")
	getCatQuery := "SELECT id, name, date_birth, vaccinated, image_path FROM cats WHERE id = $1"
	var cat model.Cat
	imageNull := sql.NullString{}
	if err := r.DB.QueryRow(ctx, getCatQuery, id).Scan(&cat.ID, &cat.Name, &cat.DateBirth, &cat.Vaccinated,
		&imageNull); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while selecting row from table cats")
		return nil, fmt.Errorf("postgres repository: can't get cat - %w", err)
	}
	if imageNull.Valid {
		cat.ImagePath = imageNull.String
	}
	return &cat, nil
}

// Update method updates object Cat from postgres database
// with selection by id and returns object Cat.
func (r CatRepository) Update(ctx context.Context, id string, input *model.UpdateCat) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Info("postgres repository: update cat")

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if input.DateBirth != nil {
		setValues = append(setValues, fmt.Sprintf("date_birth=$%d", argID))
		args = append(args, *input.DateBirth)
		argID++
	}

	if input.Vaccinated != nil {
		setValues = append(setValues, fmt.Sprintf("vaccinated=$%d", argID))
		args = append(args, *input.Vaccinated)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	updateCatQuery := fmt.Sprintf("UPDATE cats SET %s WHERE id = $%d RETURNING id, name, date_birth,"+
		" vaccinated, image_path;", setQuery, argID)
	args = append(args, id)

	var cat model.Cat
	imageNull := sql.NullString{}

	if err := r.DB.QueryRow(ctx, updateCatQuery, args...).Scan(&cat.ID, &cat.Name, &cat.DateBirth, &cat.Vaccinated,
		&imageNull); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while updating row from table cats")
		return nil, fmt.Errorf("postgres repository: can't update cat - %w", err)
	}
	if imageNull.Valid {
		cat.ImagePath = imageNull.String
	}

	return &cat, nil
}

// Delete method deletes object Cat from postgres database
// with selection by id.
func (r CatRepository) Delete(ctx context.Context, id string) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Info("repository: delete cat")
	deleteCatQuery := "DELETE FROM cats WHERE id = $1"
	_, err := r.DB.Exec(ctx, deleteCatQuery, id)
	if err != nil {
		logrus.Error(err, "postgres repository: Error occurred while deleting row from table cats")
		return fmt.Errorf("postgres repository: can't delete cat - %w", err)
	}

	return nil
}

// UploadImage method updates image path object Cat from postgres database
// with selection by id and returns object Cat.
func (r CatRepository) UploadImage(ctx context.Context, id string, path string) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Info("postgres repository: update cats image path")
	UpdateImagePathCatQuery := "UPDATE cats SET image_path=$1 WHERE id = $2 RETURNING id, name, date_birth," +
		" vaccinated, image_path;"
	var cat model.Cat
	imageNull := sql.NullString{}
	if err := r.DB.QueryRow(ctx, UpdateImagePathCatQuery, path, id).Scan(&cat.ID, &cat.Name, &cat.DateBirth,
		&cat.Vaccinated, &imageNull); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while updating image path table cats")
		return nil, fmt.Errorf("postgres repository: can't update cats image path - %w", err)
	}
	if imageNull.Valid {
		cat.ImagePath = imageNull.String
	}

	return &cat, nil
}
