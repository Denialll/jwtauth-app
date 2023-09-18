package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const (
	salt = "hjqrasfasfajfhajs"
)

type AuthService struct {
	repo            repository.Authorization
	tokenManager    pkg.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization, tokenManager pkg.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateTokens(ctx context.Context, username, password string) (Tokens, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) createSession(ctx context.Context, studentId int) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(studentId))
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken(strconv.Itoa(studentId))
	if err != nil {
		return res, err
	}

	err = s.repo.SetSession(ctx, studentId, res.RefreshToken)

	return res, err
}
