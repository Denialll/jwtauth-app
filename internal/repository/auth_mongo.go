package repository

import (
	"context"
	"github.com/Denialll/jwtauth-app/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	db *mongo.Collection
}

func NewAuthRepo(db *mongo.Database) *AuthRepo {
	return &AuthRepo{
		db: db.Collection("users"),
	}
}

func (a *AuthRepo) CreateUser(ctx context.Context, user model.User) (string, error) {
	user.Id = uuid.New().String()
	_, err := a.db.InsertOne(ctx, user)

	return user.Id, err
}

func (a *AuthRepo) GetUser(ctx context.Context, uuid string) (model.User, error) {
	var user model.User

	filter := bson.M{"_id": uuid}
	err := a.db.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (a *AuthRepo) GetUserByRefToken(ctx context.Context, refreshToken string) (model.User, error) {
	var user model.User

	filter := bson.M{"session.refreshToken": refreshToken}
	err := a.db.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (a *AuthRepo) SetSession(ctx context.Context, uuid string, session model.Session) error {
	_, err := a.db.UpdateOne(ctx, bson.M{"_id": uuid}, bson.M{"$set": bson.M{"session": session}})

	return err
}
