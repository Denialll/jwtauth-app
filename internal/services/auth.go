package services

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const (
	salt       = "hjqrasfasfajfhajs"
	signingKey = "aaaaa"
	tokenTTL   = 12 * time.Hour
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

func NewAuthService(repo repository.Authorization, tokenManager pkg.TokenManager, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
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

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
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

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.Itoa(studentId), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(tokenTTL),
	}

	err = s.repo.SetSession(ctx, studentId, session)

	return res, err
}
