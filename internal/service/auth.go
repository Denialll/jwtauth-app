package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/model"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type AuthService struct {
	repo            repository.Authorization
	tokenManager    pkg.TokenManager
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.Authorization, tokenManager pkg.TokenManager, refreshToken time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		tokenManager:    tokenManager,
		refreshTokenTTL: refreshToken,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (string, error) {
	var err error

	user.Password, err = algoBCrypt(user.Password)
	if err != nil {
		return "", err
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) GenerateTokens(ctx context.Context, uuid string) (Tokens, error) {
	user, err := s.repo.GetUser(ctx, uuid)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(ctx, user.Id)
}

func (s *AuthService) UpdateTokens(ctx context.Context, authHeader, refreshToken string) (Tokens, error) {
	headerParts := strings.Split(authHeader, " ")
	tokenClaims, err := s.tokenManager.ParseUnverified(headerParts[1])

	user, err := s.repo.GetUser(ctx, tokenClaims["sub"].(string))
	if err != nil {
		return Tokens{}, fmt.Errorf("mongoDB Exception: %s", err)
	}

	if user.Session.ExpiresAt.Unix() < time.Now().Unix() {
		return Tokens{}, fmt.Errorf("refresh token expired")
	}

	if tokenClaims["nbf"].(float64) != float64(user.Session.IssuedAt.Unix()) {
		return Tokens{}, fmt.Errorf("these tokens were not issued together")
	}

	decodedToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Session.RefreshToken), decodedToken)
	if err != nil {
		return Tokens{}, err
	}

	tokens, err := s.GenerateTokens(ctx, user.Id)
	if err != nil {
		return Tokens{}, fmt.Errorf("mongoDB Exception: %s", err)
	}

	return tokens, err
}

func (s *AuthService) createSession(ctx context.Context, uuid string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	timeCreate := time.Now()

	res.AccessToken, err = s.tokenManager.NewAccessToken(uuid, timeCreate)
	if err != nil {
		return Tokens{}, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := algoBCrypt(res.RefreshToken)
	if err != nil {
		return Tokens{}, err
	}

	session := model.Session{
		RefreshToken: refreshToken,
		ExpiresAt:    timeCreate.Add(s.refreshTokenTTL),
		IssuedAt:     timeCreate,
	}

	err = s.repo.SetSession(ctx, uuid, session)
	if err != nil {
		return Tokens{}, err
	}

	return res, err
}

func algoBCrypt(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}
