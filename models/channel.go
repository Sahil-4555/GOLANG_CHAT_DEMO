package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ChannelCollection = "channels"

var Channel_Types = [...]string{"one-to-one", "private"}

type Channel struct {
	Id                primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	ChannelName       string               `json:"channel_name" bson:"channel_name" validate:"required,max=50"`
	Slug              string               `json:"slug,omitempty" bson:"slug,omitempty"`
	Description       string               `json:"description,omitempty" bson:"description,omitempty"`
	Users             []primitive.ObjectID `json:"users" bson:"users"`
	Admins            []primitive.ObjectID `json:"admins,omitempty" bson:"admins,omitempty"`
	CloseConversation []primitive.ObjectID `json:"close_conversation,omitempty" bson:"close_conversation,omitempty"`
	ChannelType       string               `json:"channel_type" bson:"channel_type" validate:"required,oneof=one-to-one private"`
	LastOpened        []LastOpenedBy       `json:"last_opened,omitempty" bson:"last_opened,omitempty"`
	LastActivity      []LastActivityBy     `json:"last_activity,omitempty" bson:"last_activity,omitempty"`
	CreatedAt         time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time            `json:"updated_at,omitempty" bson:"upadted_at,omitempty"`
	DeletedAt         *(time.Time)         `json:"deleted_at" bson:"deleted_at"`
}

type LastOpenedBy struct {
	UserId       primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	LastOpenedAt time.Time          `json:"last_opened_at,omitempty" bson:"last_opened_at,omitempty"`
}

type LastActivityBy struct {
	UserId         primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
	LastActivityAt *(time.Time)       `json:"last_activity_at" bson:"last_activity_at"`
}

func (c *Channel) TimeStamp() {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (l *LastOpenedBy) TimeStamp() {
	l.LastOpenedAt = time.Now()
}

func (l *LastActivityBy) TimeStamp() {
	time := time.Now()
	l.LastActivityAt = &time
}

func (c *Channel) NewID() {
	c.Id = primitive.NewObjectID()
}
