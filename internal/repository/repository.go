package repository

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (string, error)
	GetUser(ctx context.Context, uuid string) (model.User, error)
	GetUserByRefToken(ctx context.Context, refreshToken string) (model.User, error)
	SetSession(ctx context.Context, uuid string, session model.Session) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthRepo(db),
	}
}
