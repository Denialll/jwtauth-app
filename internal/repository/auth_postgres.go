package repository

import (
	"context"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *AuthPostgres) CreateUser(ctx context.Context, user models.User) (string, error) {
	user.Id = uuid.New().String()
	_, err := r.db.InsertOne(ctx, user)
	fmt.Println(user.Id)

	return user.Id, err
}

func (r *AuthPostgres) GetUser(ctx context.Context, uuid string) (models.User, error) {
	fmt.Println("GUID: " + uuid)
	var user models.User
	filter := bson.M{"_id": uuid}
	err := r.db.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (r *AuthPostgres) SetSession(ctx context.Context, uuid string, refreshToken string) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": uuid}, bson.M{"$set": bson.M{"refreshToken": refreshToken}})

	return err
}
