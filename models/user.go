package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UserCollection = "user"

type User struct {
	Id         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name" validate:"required,max=50"`
	UserName   string               `json:"user_name" bson:"user_name" validate:"required,max=50"`
	Email      string               `json:"email" bson:"email" validate:"required"`
	Password   string               `json:"password" bson:"password" validate:"required"`
	Status     int                  `json:"status" bson:"status"`
	Favourites []primitive.ObjectID `json:"favourites,omitempty" bson:"favourites,omitempty"`
	CreatedAt  time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at" bson:"updated_at"`
	DeletedAt  *(time.Time)         `json:"deleted_at" bson:"deleted_at"`
}

func (u *User) TimeStamp() {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) NewUser() {
	u.Id = primitive.NewObjectID()
}
