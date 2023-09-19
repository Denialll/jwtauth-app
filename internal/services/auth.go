package services

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

//type tokenClaims struct {
//	jwt.StandardClaims
//	UserId int `json:"user_id"`
//}

func NewAuthService(repo repository.Authorization, tokenManager pkg.TokenManager) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error) {
	user.Password = generatePasswordHash(user.Password)

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateTokens(ctx context.Context, username, password string) (Tokens, error) {
	user, err := s.repo.GetUser(ctx, username, generatePasswordHash(password))
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

func (s *AuthService) createSession(ctx context.Context, userId primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId.Hex())
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken(userId.Hex())
	if err != nil {
		return res, err
	}

	err = s.repo.SetSession(ctx, userId, res.RefreshToken)

	return res, err
}
