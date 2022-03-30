package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

// AuthRepository type represents postgres behavior for authentication.
type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

// CreateUser method create user in postgres database.
func (r AuthRepository) CreateUser(ctx context.Context, input *CreateUserInput) error {
	logrus.WithFields(logrus.Fields{
		"ID":           input.ID,
		"userName":     input.UserName,
		"email":        input.Email,
		"password":     input.Password,
		"refreshToken": input.RefreshToken,
	}).Info("postgres repository: create User")
	_, err := r.DB.Exec(ctx, `INSERT INTO USERS (id, name, email, password, refresh_token) VALUES($1, $2, $3, $4, $5)`,
		input.ID, input.UserName, input.Email, input.Password, input.RefreshToken)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key"):
			switch {
			case strings.Contains(err.Error(), "users_email_key"):
				logrus.Error(err, "postgres repository: can't create User")
				return errors.New("user with given email exists, change your email")

			case strings.Contains(err.Error(), "users_primary_key"):
				logrus.Error(err, "postgres repository: can't create User")
				return errors.New("user with given UUID exists, try to create again")
			}

		default:
			logrus.Error(err, "postgres repository: can't create User")
			return errors.New("can't create User")
		}
	}

	return nil
}

// GetUserHashedPassword method returns user id with hashed password from postgres database.
func (r AuthRepository) GetUserHashedPassword(ctx context.Context, email string) (string, string, error) {
	logrus.WithFields(logrus.Fields{
		"email": email,
	}).Info("postgres repository: get user id and hashed password")
	var id, password string
	if err := r.DB.QueryRow(ctx, `SELECT id, password FROM users WHERE email = $1`,
		email).Scan(&id, &password); err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			logrus.Error(err, "postgres repository: user with given email doesn't exist")
			return "", "", errors.New("user with given email doesn't exist")

		default:
			logrus.Error(err, "postgres repository: can't get users hashed password")
			return "", "", errors.New("can't get users hashed password")
		}
	}

	return id, password, nil
}

// GetUserRefreshToken method returns users stored refresh token.
func (r AuthRepository) GetUserRefreshToken(ctx context.Context, id string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"id": id,
	}).Info("postgres repository: get user refresh token by his id")
	var refreshTokenString string
	if err := r.DB.QueryRow(ctx, `SELECT refresh_token FROM users WHERE id = $1`,
		id).Scan(&refreshTokenString); err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			logrus.Error(err, "postgres repository: user with given UUID doesn't exist")
			return "", errors.New("user with given UUID doesn't exist")

		default:
			logrus.Error(err, "postgres repository: can't get users refresh token")
			return "", errors.New("can't get users refresh token")
		}
	}

	return refreshTokenString, nil
}

// UpdateUserRefreshToken method updates user refresh token.
func (r AuthRepository) UpdateUserRefreshToken(ctx context.Context, id, refreshToken string) error {
	logrus.WithFields(logrus.Fields{
		"id":           id,
		"refreshToken": refreshToken,
	}).Info("postgres repository: update refresh token")
	if _, err := r.DB.Exec(ctx, `UPDATE users SET refresh_token = $1 WHERE id = $2`,
		refreshToken, id); err != nil {
		switch {
		case strings.Contains(err.Error(), "no rows in result set"):
			logrus.Error(err, "postgres repository: user with given UUID doesn't exist")
			return errors.New("user with given UUID doesn't exist")

		default:
			logrus.Error(err, "postgres repository: can't update refresh token")
			return errors.New("can't update refresh token")
		}
	}

	return nil
}
