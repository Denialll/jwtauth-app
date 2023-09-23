package model

type User struct {
	Id       string  `json:"-" bson:"_id,omitempty"`
	Email    string  `json:"email" bson:"email" binding:"required"`
	Password string  `json:"password" bson:"password" binding:"required"`
	Session  Session `json:"session" bson:"session,omitempty"`
}
