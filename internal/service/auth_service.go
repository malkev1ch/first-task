package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/malkev1ch/first-task/internal/repository"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExTime  = 15
	refreshTokenExTime = 720
)

type SignUpInput struct {
	// The Name of a user
	// example: Some name
	// required: true
	UserName string `json:"userName"`
	// The email of a user
	// example: qwerty@gmail.com
	// required: true
	Email string `json:"email"`
	// The password of a user
	// example: ZAQ!2wsx
	// required: true
	Password string `json:"password"`
}

type SignInInput struct {
	// The email of a user
	// example: qwerty@gmail.com
	// required: true
	Email string `json:"email"`
	// The password of a user
	// example: ZAQ!2wsx
	// required: true
	Password string `json:"password"`
}

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// SignUp method hash user password and after that save user in repository.
func (s Service) SignUp(ctx context.Context, input *SignUpInput) (*model.Tokens, error) {
	hPassword, err := s.hashPassword(input.Password)
	if err != nil {
		logrus.Error(err, "service: hash password failed")
		return nil, fmt.Errorf("service: hash password failed - %w", err)
	}

	input.Password = hPassword
	id := uuid.New().String()
	tokens, err := s.generateToken(input.Email, id)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(ctx, &repository.CreateUserInput{
		ID: id, UserName: input.UserName, Email: input.Email,
		Password: input.Password, RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// SignIn Generates tokens for created user.
func (s Service) SignIn(ctx context.Context, input *SignInInput) (*model.Tokens, error) {
	id, hash, err := s.repo.GetUserHashedPassword(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if pass := s.checkPasswordHash(input.Password, hash); !pass {
		logrus.Error(err, "service: incorrect password")
		return nil, fmt.Errorf("service: incorrect password - %w", err)
	}

	tokens, err := s.generateToken(input.Email, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUserRefreshToken(ctx, id, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

// RefreshToken method checks refresh token for validity and if it's ok return new token pair.
func (s Service) RefreshToken(ctx context.Context, refreshTokenString string) (*model.Tokens, error) {
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		logrus.Error(err, "service: can't parse refresh token")
		return nil, fmt.Errorf("service: can't parse refresh token - %w", err)
	}
	if !refreshToken.Valid {
		logrus.Info("service: expired refresh token")
		return nil, fmt.Errorf("service: expired refresh token")
	}
	claims := refreshToken.Claims.(jwt.MapClaims)
	userID := claims["jti"]
	email := claims["email"]
	fmt.Println(userID, email)
	if userID == "" || email == "" {
		logrus.Error(err, "service: error while parsing claims", userID, email)
		return nil, fmt.Errorf("service: error while parsing claims")
	}

	lastRefreshTokenStringLast, err := s.repo.GetUserRefreshToken(ctx, userID.(string))
	if err != nil {
		return nil, fmt.Errorf("service: token refresh failed - %w", err)
	}
	if refreshTokenString != lastRefreshTokenStringLast {
		logrus.Error(err, "service: invalid refresh token")
		return nil, fmt.Errorf("service: invalid refresh token")
	}

	tokens, err := s.generateToken(fmt.Sprintf("%v", email), fmt.Sprintf("%v", userID))
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateUserRefreshToken(ctx, fmt.Sprintf("%v", userID), tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

// hashPassword from string
// bcrypt.DefaultCost = 10.
func (s Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compare encrypt.
func (s Service) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken generates new token pair and with putting email inside payload.
func (s Service) generateToken(email string, id string) (*model.Tokens, error) {
	expirationTimeAT := time.Now().Add(accessTokenExTime * time.Minute)
	expirationTimeRT := time.Now().Add(refreshTokenExTime * time.Hour)

	accessTokensClaims := &JwtCustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			ExpiresAt: expirationTimeAT.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokensClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		logrus.Error(err, "service: can't generate access token")
		return nil, fmt.Errorf("service: can't generate access token - %w", err)
	}

	refreshTokensClaims := &JwtCustomClaims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			Id:        id,
			ExpiresAt: expirationTimeRT.Unix(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokensClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		logrus.Error(err, "service: can't generate refresh token")
		return nil, fmt.Errorf("service: can't generate refresh token - %w", err)
	}

	return &model.Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}
