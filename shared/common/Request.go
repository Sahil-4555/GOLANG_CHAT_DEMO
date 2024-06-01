package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model SignUpRequest
type SignUpReq struct {
	// Username of the user
	//
	// required: true
	// example: john_doe
	UserName string `json:"user_name" bson:"user_name" validate:"required,max=50"`

	// Name of the user
	//
	// required: true
	// example: john_doe
	Name string `json:"name" bson:"name" validate:"required,max=50"`

	// Email of the user
	//
	// required: true
	// example: john@example.com
	Email string `json:"email" bson:"email" validate:"required"`

	// Password of the user
	//
	// required: true
	// example: password123
	Password string `json:"password" bson:"password" validate:"required"`
}

// swagger:model SignInRequest
type SignInReq struct {
	// Email of the user
	//
	// required: true
	// example: john@example.com
	Email string `json:"email" bson:"email" validate:"required"`

	// Password of the user
	//
	// required: true
	// example: password123
	Password string `json:"password" bson:"password" validate:"required"`
}

type RoomReq struct {
	Sender  primitive.ObjectID `json:"sender" bson:"sender"`
	Reciver primitive.ObjectID `json:"reciver" bson:"reciver"`
}

type JoinChannelReq struct {
	Reciver primitive.ObjectID `json:"reciver_id" bson:"reciver_id"`
}

type CreateGroupReq struct {
	ChannelName string `json:"channel_name" bson:"channel_name" validate:"required,max=50"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
}

type AddMessagesReq struct {
	Sender      string             `json:"sender" bson:"sender"`
	ChannelId   primitive.ObjectID `json:"channel_id" bson:"channel_id"`
	Content     string             `json:"content" bson:"content"`
	ContentType string             `json:"content_type" bson:"content_type"`
	MediaUrl    string             `json:"media_url" bson:"media_url"`
}

type NewMessageReq struct {
	Sender      string             `json:"sender" bson:"sender" validate:"required"`
	ChannelId   string             `json:"channel_id" bson:"channel_id" validate:"required"`
	ParentId    primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Content     string             `json:"content" bson:"content"`
	ContentType string             `json:"content_type" bson:"content_type" validate:"required,oneof=text media both"`
	MediaUrl    string             `json:"media_url,omitempty" bson:"media_url,omitempty"`
}

type AddMembersToGroupReq struct {
	ChannelId primitive.ObjectID   `json:"channel_id" bson:"channel_id" validate:"required"`
	Users     []primitive.ObjectID `json:"users" bson:"users" validate:"required"`
}

type PostNotificationOnMentionReq struct {
	Users  []string `json:"users" bson:"users"`
	Sender string   `json:"sender" bson:"sender"`
}

type JoinRoomReq struct {
	ChannelId string             `json:"channel_id" bson:"channel_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type GiveAdminRightsToUserReq struct {
	ChannelId string             `json:"channel_id" bson:"channel_id"`
	UserId    primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type AddFavoriteChannelReq struct {
	IsFavourite bool               `json:"is_favourite" bson:"is_favourite"`
	ChannelID   primitive.ObjectID `json:"channel_id" bson:"channel_id"`
}

type UpdateMessageReq struct {
	Id      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content,omitempty" bson:"content,omitempty"`
}

type DeleteMessageReq struct {
	Id primitive.ObjectID `json:"_id" bson:"_id"`
}

type StatusChangeReq struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}

type LeaveChannelReq struct {
	ChannelId primitive.ObjectID `json:"channel_id" bson:"channel_id" validate:"required"`
}

type RemovUserFromGroupByGroupAdminReq struct {
	ChannelId primitive.ObjectID   `json:"channel_id" bson:"channel_id" validate:"required"`
	Users     []primitive.ObjectID `json:"users" bson:"users" validate:"required"`
}

type CloseConversationReq struct {
	ChannelId primitive.ObjectID `json:"channel_id" bson:"channel_id" validate:"required"`
}

type UploadMediaReq struct {
	Key      string             `bson:"key" json:"key" validate:"required"`
	FileName string             `bson:"file_name" json:"file_name" validate:"required"`
	UserId   primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	FileData []byte             `bson:"file_data" json:"file_data" validate:"required"`
	FileSize int64              `bson:"file_size" json:"file_size" validate:"required"`
}
