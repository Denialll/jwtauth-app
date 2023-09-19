package services

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/Denialll/jwtauth-app/internal/repository"
	"github.com/Denialll/jwtauth-app/pkg"
	"time"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GenerateTokens(ctx context.Context, uuid string) (Tokens, error)
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TodoList interface {
}

type TodoItem interface {
}

type Deps struct {
	Repos                 *repository.Repository
	TokenManager          pkg.TokenManager
	AccessTTL, RefreshTTL time.Duration
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(deps Deps) *Service {
	return &Service{
		Authorization: NewAuthService(deps.Repos, deps.TokenManager),
	}
}
