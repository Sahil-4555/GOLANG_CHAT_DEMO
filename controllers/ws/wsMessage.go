package ws

import (
	"encoding/json"

	"time"

	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MAX_ROOM_CAPACITY              = 50 // Max Room Capacity
	SendMessageAction              = "send-message"
	JoinRoomAction                 = "join-room"
	LeaveRoomAction                = "leave-room"
	ClientNameAction               = "client-name"
	JoinRoomNotification           = "join-room-notify"
	ClientListNotification         = "client-list-notify"
	FailJoinRoomNotification       = "fail-join-room-notify"
	KnowYourUser                   = "know-yourself" // Individual message with client details send to a particulr user after he enters a room contains slug and name
	PostGlobalNotification         = "post-notification-to-user"
	HandleMessage                  = "handle-message"
	GetGlobalNotification          = "get-notification"
	Error                          = "error"
	ChangeUserStatus               = "change-user-status"
	UserStatusOnline               = "user-status-online"
	UserStatusAway                 = "user-status-away"
	UserStatusOffline              = "user-status-offline"
	UserStatusDoNotDisturb         = "user-status-do-not-disturb"
	UpdateChannelDataAcrossUser    = "update-channel-data-across-user"
	MetionInMessageNotification    = "mention-in-message-notification"
	UpdateChannelDataAcrossChannel = "update-channel-data-across-channel"
	UpdateChannelOnMessage         = "update-channel-on-message"
	UpdateMessage                  = "update-message"
	DeleteMessage                  = "delete-message"
	UserTyping                     = "user-typing"
	UserStopTyping                 = "user-stop-typing"
	AddChannelOnAddingMember       = "add-channel-on-adding-member"
)

type ClientsinRoomMessage struct { // we are using this to return list of clients to all clients in room when register unregister happens
	Action     string    `json:"action"`     //action
	ClientList []*Client `json:"clientlist"` //message
	Target     string    `json:"target"`     //target the room
	Sender     *Client   `json:"sender"`     //whose readpump is used
}

type GlobalNotificationToUser struct {
	Action string  `json:"action"`
	Target string  `json:"target"`
	Sender *Client `json:"sender"`
}

type Message struct { // in readpump also can be used in writepump
	Action  string  `json:"action"`            // action
	Message string  `json:"message,omitempty"` // message
	Target  string  `json:"target,omitempty"`  // target the room
	Sender  *Client `json:"sender,omitempty"`  // whose readpump is used

	/* This are the field for the send message */
	ContentType int                  `json:"content_type,omitempty"` // content type for the message
	NotifyUsers []primitive.ObjectID `json:"notify_users,omitempty"` // userids which you want to notify.
	ParentId    primitive.ObjectID   `json:"parent_id,omitempty"`    // parentid of the message if there for reply to the message
	FileName    string               `json:"file_name,omitempty"`    // media_url for the media sharing

	/* This are the fields for the channel data and related notification sharing across the channel user joined */
	NotificationType string                   `json:"notification_type,omitempty"`
	Channel          []common.ChannelResponse `json:"channel,omitempty"`
	ChannelId        string                   `json:"channel_id,omitempty"`
	Id               primitive.ObjectID       `json:"_id,omitempty"`
	UpdatedAt        time.Time                `json:"updated_at,omitempty"`
	CreatedAt        time.Time                `json:"created_at,omitempty"`
	Payload          interface{}              `json:"payload,omitempty"`

	/* This are the fields for the changing user status */
	Status int `json:"status,omitempty"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
	}

	return json
}

func (message *ClientsinRoomMessage) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
	}

	return json
}
