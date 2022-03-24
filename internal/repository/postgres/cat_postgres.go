package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

// Create method saves object Cat into postgres database.
func (r RepositoryPostgres) Create(ctx context.Context, input *model.CreateCat) (string, error) {
	id := uuid.New().String()
	logrus.WithFields(logrus.Fields{
		"id":         id,
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("postgres repository: create cat")

	insertCatQuery := "INSERT INTO cats(id, name, date_birth, vaccinated) VALUES ($1, $2, $3, $4)"

	if _, err := r.DB.Exec(ctx, insertCatQuery, id, input.Name, input.DateBirth, input.Vaccinated); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while inserting new row in table cats")
		return "", fmt.Errorf("postgres repository: can't create cat - %w", err)
	}
	return id, nil
}

// Get method returns object Cat from postgres database
// with selection by id.
func (r RepositoryPostgres) Get(ctx context.Context, id string) (model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("postgres repository: get cat")
	getCatQuery := "SELECT id, name, date_birth, vaccinated, image_path FROM cats WHERE id = $1"
	var cat model.Cat
	if err := r.DB.QueryRow(ctx, getCatQuery, id).Scan(&cat.ID, &cat.Name, &cat.DateBirth, &cat.Vaccinated,
		&cat.ImagePath); err != nil {
		logrus.Error(err, "postgres repository: Error occurred while selecting row from table cats")
		return cat, fmt.Errorf("postgres repository: can't get cat - %w", err)
	}

	return cat, nil
}

// Update method updates object Cat from postgres database
// with selection by id.
func (r RepositoryPostgres) Update(ctx context.Context, id string, input *model.UpdateCat) error {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("postgres repository: update cat")

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
	updateCatQuery := fmt.Sprintf("UPDATE cats SET %s WHERE id = $%d", setQuery, argID)
	args = append(args, id)

	_, err := r.DB.Exec(ctx, updateCatQuery, args...)
	if err != nil {
		logrus.Error(err, "postgres repository: Error occurred while updating row from table cats")
		return fmt.Errorf("postgres repository: can't update cat - %w", err)
	}

	return nil
}

// Delete method deletes object Cat from postgres database
// with selection by id.
func (r RepositoryPostgres) Delete(ctx context.Context, id string) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("repository: delete cat")
	deleteCatQuery := "DELETE FROM cats WHERE id = $1"
	_, err := r.DB.Exec(ctx, deleteCatQuery, id)
	if err != nil {
		logrus.Error(err, "postgres repository: Error occurred while deleting row from table cats")
		return fmt.Errorf("postgres repository: can't delete cat - %w", err)
	}

	return nil
}

// UploadImage method updates image path object Cat from postgres database
// with selection by id.
func (r RepositoryPostgres) UploadImage(ctx context.Context, id, path string) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("postgres repository: update cats image path")
	UpdateImagePathCatQuery := "UPDATE cats SET image_path=$1 WHERE id = $2"
	_, err := r.DB.Exec(ctx, UpdateImagePathCatQuery, path, id)
	if err != nil {
		logrus.Error(err, "postgres repository: Error occurred while updating image path table cats")
		return fmt.Errorf("postgres repository: can't update cats image path - %w", err)
	}

	return nil
}
