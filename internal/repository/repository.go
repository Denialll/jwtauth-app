package repository

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (string, error)
	GetUser(ctx context.Context, uuid string) (models.User, error)
	SetSession(ctx context.Context, uuid string, refreshToken string) error
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
