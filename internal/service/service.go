package service

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/model"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"time"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GenerateTokens(ctx context.Context, uuid string) (Tokens, error)
	UpdateTokens(ctx context.Context, accessToken, refreshToken string) (Tokens, error)
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Deps struct {
	Repos                 *repository.Repository
	TokenManager          pkg.TokenManager
	AccessTTL, RefreshTTL time.Duration
}

type Service struct {
	Authorization
}

func NewService(deps Deps) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repos, deps.TokenManager, deps.RefreshTTL),
	}
}
