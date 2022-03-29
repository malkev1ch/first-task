package repository

import (
	"context"
	"fmt"

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
		logrus.Error(err, "postgres repository: can't create User")
		return fmt.Errorf("postgres repository: can't create User - %w", err)
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
		logrus.Error(err, "postgres repository: can't get user hashed password")
		return "", "", fmt.Errorf("postgres repository: can't get user hashed password - %w", err)
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
		logrus.Error(err, "postgres repository: can't get user refresh token")
		return "", fmt.Errorf("postgres repository: can't get user refresh token - %w", err)
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
		logrus.Error(err, "postgres repository: can't update refresh token")
		return fmt.Errorf("postgres repository: can't update refresh token - %w", err)
	}
	return nil
}
