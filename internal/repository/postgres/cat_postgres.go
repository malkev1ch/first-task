package postgres

import (
	"context"
	"fmt"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
	"strings"
)

func (r RepositoryPostgres) Create(ctx context.Context, input *model.Cat) (int, error) {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("repository: create cat")
	insertCatQuery := fmt.Sprintf("INSERT INTO cats(name, date_birth, vaccinated) VALUES ($1, $2, $3) RETURNING row_id")
	var id int

	if err := r.DB.QueryRow(ctx, insertCatQuery, input.Name, input.DateBirth, input.Vaccinated).Scan(&id); err != nil {
		logrus.Error(err, "Error occurred while inserting new row in table cats")
		return 0, err
	}
	return id, nil
}

func (r RepositoryPostgres) Get(ctx context.Context, id int) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("repository: get cat")
	getCatQuery := fmt.Sprintf("SELECT name, date_birth, vaccinated, image_path FROM cats WHERE row_id = $1")
	var cat model.Cat
	if err := r.DB.QueryRow(ctx, getCatQuery, id).Scan(&cat.Name, &cat.DateBirth, &cat.Vaccinated, &cat.ImagePath); err != nil {
		logrus.Error(err, "Error occurred while selecting row from table cats")
		return nil, err
	}

	return &cat, nil
}

func (r RepositoryPostgres) Update(ctx context.Context, id int, input *model.Cat) error {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("repository: update cat")

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
	updateCatQuery := fmt.Sprintf("UPDATE cats SET %s WHERE row_id = $%d", setQuery, argId)
	args = append(args, id)

	_, err := r.DB.Exec(ctx, updateCatQuery, args...)
	if err != nil {
		logrus.Error(err, "Error occurred while updating row from table cats")
		return err
	}

	return nil
}

func (r RepositoryPostgres) Delete(ctx context.Context, id int) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("repository: delete cat")
	deleteCatQuery := fmt.Sprintf("DELETE FROM cats WHERE row_id = $1")
	_, err := r.DB.Exec(ctx, deleteCatQuery, id)
	if err != nil {
		logrus.Error(err, "Error occurred while deleting row from table cats")
		return err
	}

	return nil
}

func (r RepositoryPostgres) UploadImage(ctx context.Context, id int, path string) error {
	UpdateImagePathCatQuery := fmt.Sprintf("UPDATE cats SET image_path=$1 WHERE row_id = $2")
	_, err := r.DB.Exec(ctx, UpdateImagePathCatQuery, path, id)
	if err != nil {
		logrus.Error(err, "Error occurred while updating image path table cats")
		return err
	}

	return nil
}
