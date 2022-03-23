package mongodb

import (
	"context"
	"fmt"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

func (r RepositoryMongo) Create(ctx context.Context, input *model.Cat) (int, error) {
	logrus.WithFields(logrus.Fields{
		"Name":       input.Name,
		"DateBirth":  input.DateBirth,
		"Vaccinated": input.Vaccinated,
	}).Debugf("repository: create cat")
	var id int

	return id, nil
}

func (r RepositoryMongo) Get(ctx context.Context, id int) (*model.Cat, error) {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("repository: get cat")
	var cat model.Cat

	return &cat, nil
}

func (r RepositoryMongo) Update(ctx context.Context, id int, input *model.Cat) error {
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

	//setQuery := strings.Join(setValues, ", ")
	args = append(args, id)

	return nil
}

func (r RepositoryMongo) Delete(ctx context.Context, id int) error {
	logrus.WithFields(logrus.Fields{
		"ID": id,
	}).Debugf("repository: delete cat")

	return nil
}

func (r RepositoryMongo) UploadImage(ctx context.Context, id int, path string) error {

	return nil
}
