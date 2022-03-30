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
	accessTokenExTime  = 720
	refreshTokenExTime = 720
)

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

// SignUp method hash user password and after that save user in repository.
func (s AuthService) SignUp(ctx context.Context, input *model.CreateUser) (*model.Tokens, error) {
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

	err = s.repo.Auth.CreateUser(ctx, &repository.CreateUserInput{
		ID: id, UserName: input.UserName, Email: input.Email,
		Password: input.Password, RefreshToken: tokens.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// SignIn Generates tokens for created user.
func (s AuthService) SignIn(ctx context.Context, input *model.AuthUser) (*model.Tokens, error) {
	id, hash, err := s.repo.Auth.GetUserHashedPassword(ctx, input.Email)
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

	if err := s.repo.Auth.UpdateUserRefreshToken(ctx, id, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

// RefreshToken method checks refresh token for validity and if it's ok return new token pair.
func (s AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*model.Tokens, error) {
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

	if err := s.repo.Auth.UpdateUserRefreshToken(ctx, fmt.Sprintf("%v", userID), tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}

// hashPassword from string
// bcrypt.DefaultCost = 10.
func (s AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compare encrypt.
func (s AuthService) checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken generates new token pair and with putting email inside payload.
func (s AuthService) generateToken(email string, id string) (*model.Tokens, error) {
	expirationTimeAT := time.Now().Add(accessTokenExTime * time.Hour)
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
