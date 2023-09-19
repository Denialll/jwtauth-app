package models

type User struct {
	Id           string `json:"-" bson:"_id,omitempty"`
	Name         string `json:"name" bson:"name" binding:"required"`
	Username     string `json:"username" bson:"username" binding:"required"`
	Password     string `json:"password" bson:"password" binding:"required"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}
