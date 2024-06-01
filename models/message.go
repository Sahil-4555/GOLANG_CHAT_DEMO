package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MessageCollection = "message"

var Content_Types = []string{"text", "media", "both"}

type Message struct {
	Id          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Sender      primitive.ObjectID   `json:"sender" bson:"sender" validate:"required"`
	ChannelId   primitive.ObjectID   `json:"channel_id" bson:"channel_id" validate:"required"`
	Content     string               `json:"content,omitempty" bson:"content,omitempty"`
	ReadBy      []primitive.ObjectID `json:"read_by,omitempty" bson:"read_by,omitempty"`
	ContentType string               `json:"content_type" bson:"content_type" validate:"required,oneof=text media both"`
	MediaUrl    string               `json:"media_url,omitempty" bson:"media_url,omitempty"`
	ParentId    primitive.ObjectID   `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	CreatedAt   time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at" bson:"updated_at"`
	DeletedAt   *(time.Time)         `json:"deleted_at" bson:"deleted_at"`
}

func (m *Message) TimeStamp() {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

func (m *Message) NewID() {
	m.Id = primitive.NewObjectID()
}
