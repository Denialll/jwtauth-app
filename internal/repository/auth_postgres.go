package repository

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthPostgres struct {
	db *mongo.Collection
}

func NewAuthPostgres(db *mongo.Database) *AuthPostgres {
	return &AuthPostgres{
		db: db.Collection("users"),
	}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error) {
	res, err := r.db.InsertOne(ctx, user)

	return res.InsertedID.(primitive.ObjectID), err
}

func (r *AuthPostgres) GetUser(ctx context.Context, username, password string) (models.User, error) {
	var user models.User
	filter := bson.M{
		"username": username,
		"password": password,
	}
	err := r.db.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (r *AuthPostgres) SetSession(ctx context.Context, userId primitive.ObjectID, refreshToken string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": bson.M{"refreshToken": refreshToken}})

	return err
}
