package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name" binding:"required"`
	Username     string             `json:"username" bson:"username" binding:"required"`
	Password     string             `json:"password" bson:"password" binding:"required"`
	RefreshToken string             `json:"refreshToken" bson:"refreshToken"`
}
