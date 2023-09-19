package repository

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error)
	GetUser(ctx context.Context, username, password string) (models.User, error)
	SetSession(ctx context.Context, userId primitive.ObjectID, refreshToken string) error
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
