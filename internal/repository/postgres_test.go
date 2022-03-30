package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var (
	db   *pgxpool.Pool
	repo *Repository
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=qwerty",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://postgres:qwerty@%s/postgres?sslmode=disable", hostAndPort)

	log.Info("Connecting to database on url: ", databaseURL)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = pgxpool.Connect(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		err = db.Ping(context.Background())
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	repo = NewRepositoryPostgres(db)

	log.Info("Created postgres repository successfully")

	flywayPath := fmt.Sprintf("-url=jdbc:postgresql://%s/postgres", hostAndPort)
	flywayConfFile := "-configFiles=/home/andreimalkevich/first-task/sql/dockertest/flyway.conf"
	cmd := exec.Command("flyway", flywayPath, flywayConfFile, "migrate")
	log.Infof("Start migration with cmd command: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command finished with error: %v", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	testTable := []struct {
		name          string
		input         *CreateUserInput
		ctx           context.Context
		expectedError error
	}{
		{
			name: "OK",
			input: &CreateUserInput{
				ID:           id,
				UserName:     "Some Name",
				Email:        "example@outlook.com",
				Password:     "qwerty",
				RefreshToken: "1234",
			},
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name: "Invalid UUID",
			input: &CreateUserInput{
				ID:           "123",
				UserName:     "Some Name",
				Email:        "example@gmail.com",
				Password:     "qwerty",
				RefreshToken: "1234",
			},
			ctx:           ctx,
			expectedError: errors.New("can't create User"),
		},
		{
			name: "User with given email exists",
			input: &CreateUserInput{
				ID:           uuid.New().String(),
				UserName:     "Some Name",
				Email:        "example@outlook.com",
				Password:     "qwerty",
				RefreshToken: "1234",
			},
			ctx:           ctx,
			expectedError: errors.New("user with given email exists, change your email"),
		},
		{
			name: "User with given UUID exists",
			input: &CreateUserInput{
				ID:           id,
				UserName:     "Some Name",
				Email:        "example@outlook.com",
				Password:     "qwerty",
				RefreshToken: "1234",
			},
			ctx:           ctx,
			expectedError: errors.New("user with given UUID exists, try to create again"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Auth.CreateUser(testCase.ctx, testCase.input)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestGetUserHashedPassword(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Auth.CreateUser(context.Background(), &CreateUserInput{
		ID:           id,
		UserName:     "Some Name",
		Email:        "TestGetUserHashedPassword@outlook.com",
		Password:     "qwerty",
		RefreshToken: "1234",
	})
	if err != nil {
		t.Fail()
	}

	testTable := []struct {
		name          string
		email         string
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			email:         "TestGetUserHashedPassword@outlook.com",
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "User doesn't exist",
			email:         "",
			ctx:           ctx,
			expectedError: errors.New("user with given email doesn't exist"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, _, err := repo.Auth.GetUserHashedPassword(testCase.ctx, testCase.email)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestGetUserRefreshToken(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Auth.CreateUser(context.Background(), &CreateUserInput{
		ID:           id,
		UserName:     "Some Name",
		Email:        "TestGetUserRefreshToken@outlook.com",
		Password:     "qwerty",
		RefreshToken: "1234",
	})
	if err != nil {
		t.Fail()
	}

	testTable := []struct {
		userID        string
		name          string
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			userID:        id,
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "User doesn't exist",
			userID:        uuid.New().String(),
			ctx:           ctx,
			expectedError: errors.New("user with given UUID doesn't exist"),
		},
		{
			name:          "Invalid UUID",
			userID:        "",
			ctx:           ctx,
			expectedError: errors.New("can't get users refresh token"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Auth.GetUserRefreshToken(testCase.ctx, testCase.userID)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUpdateUserRefreshToken(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Auth.CreateUser(context.Background(), &CreateUserInput{
		ID:           id,
		UserName:     "Some Name",
		Email:        "TestUpdateUserRefreshToken@outlook.com",
		Password:     "qwerty",
		RefreshToken: "1234",
	})
	if err != nil {
		t.Fail()
	}

	testTable := []struct {
		userID        string
		refreshToken  string
		name          string
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			userID:        id,
			refreshToken:  "",
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "Invalid UUID",
			userID:        "",
			refreshToken:  "",
			ctx:           ctx,
			expectedError: errors.New("can't update refresh token"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Auth.UpdateUserRefreshToken(testCase.ctx, testCase.userID, testCase.refreshToken)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestCreateCat(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	testTable := []struct {
		name          string
		input         *model.Cat
		ctx           context.Context
		expectedError error
	}{
		{
			name: "OK",
			input: &model.Cat{
				ID: id,
			},
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name: "User with given email exists",
			input: &model.Cat{
				ID: id,
			},
			ctx:           ctx,
			expectedError: errors.New("cat with given UUID already exists, try to create again"),
		},
		{
			name:          "Invalid UUID",
			input:         &model.Cat{},
			ctx:           ctx,
			expectedError: errors.New("can't create cat"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Cat.Create(testCase.ctx, testCase.input)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestGetCat(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Cat.Create(context.Background(), &model.Cat{
		ID: id,
	})
	if err != nil {
		t.Fail()
	}
	testTable := []struct {
		name          string
		input         string
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			input:         id,
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "Cat with given UUID doesn't exist",
			input:         uuid.New().String(),
			ctx:           ctx,
			expectedError: errors.New("cat with given UUID doesn't exist"),
		},
		{
			name:          "Invalid UUID",
			input:         "123",
			ctx:           ctx,
			expectedError: errors.New("can't get cat"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Cat.Get(testCase.ctx, testCase.input)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestUpdateCat(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Cat.Create(context.Background(), &model.Cat{
		ID: id,
	})
	if err != nil {
		t.Fail()
	}

	name := "some Name"
	dateBirth := time.Now()
	vaccinated := false

	testUpdateCat := &model.UpdateCat{
		Name:       &name,
		DateBirth:  &dateBirth,
		Vaccinated: &vaccinated,
	}

	testTable := []struct {
		name          string
		catID         string
		input         *model.UpdateCat
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			catID:         id,
			input:         testUpdateCat,
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "Cat with given UUID doesn't exist",
			catID:         uuid.New().String(),
			input:         testUpdateCat,
			ctx:           ctx,
			expectedError: errors.New("cat with given UUID doesn't exists"),
		},
		{
			name:          "Invalid UUID",
			catID:         uuid.New().String(),
			input:         testUpdateCat,
			ctx:           ctx,
			expectedError: errors.New("cat with given UUID doesn't exists"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Cat.Update(testCase.ctx, testCase.catID, testCase.input)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}

func TestDeleteCat(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	err := repo.Cat.Create(context.Background(), &model.Cat{
		ID: id,
	})
	if err != nil {
		t.Fail()
	}
	testTable := []struct {
		name          string
		catID         string
		ctx           context.Context
		expectedError error
	}{
		{
			name:          "OK",
			catID:         id,
			ctx:           ctx,
			expectedError: nil,
		},
		{
			name:          "Invalid UUID",
			catID:         "",
			ctx:           ctx,
			expectedError: errors.New("can't delete cat"),
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Cat.Delete(testCase.ctx, testCase.catID)
			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
