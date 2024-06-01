package common

import (
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchHandlerResposne struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" structs:"_id"`
	ChannelName string             `json:"channel_name" bson:"channel_name"  structs:"channel_name"`
	ChannelType string             `json:"channel_type" bson:"channel_type" structs:"channel_type"`
	LastOpened  []time.Time        `json:"last_opened" bson:"last_opened" structs:"last_opened"`
}

type AllUserResponse struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" structs:"_id"`
	Name     string             `json:"name" bson:"name"  structs:"name"`
	Email    string             `json:"email" bson:"email" structs:"email"`
	UserName string             `json:"user_name" bson:"user_name" structs:"user_name"`
}

type ResponseData struct {
	Data interface{} `json:"data" structs:"data"`
	Meta interface{} `json:"meta" structs:"meta"`
}

type BaseSuccessResponse struct {
	Data interface{} `json:"data" bson:"data"  structs:"data"`
	Meta Meta        `json:"meta" bson:"meta"  structs:"meta"`
}

type Meta struct {
	Message string `json:"message" structs:"message"`
	Code    int    `json:"code" structs:"code"`
	Token   string `json:"token,omitempty" structs:"token,omitempty" bson:"token,omitempty"`
}

type LoginResponse struct {
	ID        primitive.ObjectID `json:"_id" structs:"_id"`
	Email     string             `json:"email" structs:"email"`
	Name      string             `json:"name" structs:"name"`
	UserName  string             `json:"user_name" structs:"user_name"`
	LastLogin time.Time          `json:"last_login" structs:"last_login"`
}

type UserResponse struct {
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email    string             `json:"email" bson:"email"`
	Status   int                `json:"status" bson:"status"`
	Name     string             `json:"name" bson:"name"`
	UserName string             `json:"user_name" bson:"user_name"`
}

type ChannelsUserResponse struct {
	UserId   primitive.ObjectID `json:"user_id" bson:"user_id"`
	Email    string             `json:"email" bson:"email"`
	Status   int                `json:"status" bson:"status"`
	Name     string             `json:"name" bson:"name"`
	UserName string             `json:"user_name" bson:"user_name"`
	IsAdmin  bool               `json:"is_admin" bson:"is_admin"`
}

type ChannelResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	ChannelType  string             `json:"channel_type" bson:"channel_type"`
	ChannelName  string             `json:"channel_name" bson:"channel_name"`
	Users        []UserResponse     `json:"users" bson:"users"`
	LastActivity []time.Time        `json:"last_activity" bson:"last_activity"`
	MessageCount int                `json:"message_count" bson:"message_count"`
}

type JoinGroupChannelResponse struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	ChannelType  string             `json:"channel_type" bson:"channel_type"`
	ChannelName  string             `json:"channel_name" bson:"channel_name"`
	Users        []UserResponse     `json:"users" bson:"users"`
	LastActivity []time.Time        `json:"last_activity" bson:"last_activity"`
	MessageCount int                `json:"message_count" bson:"message_count"`
	IsFavourite  bool               `json:"is_favourites" bson:"is_favourites"`
}

type GetChannelMembersReposne struct {
	ID    primitive.ObjectID     `json:"_id" bson:"_id"`
	Users []ChannelsUserResponse `json:"users" bson:"users"`
}

type GetOnetoOneFavouriteChannelResponse struct {
	ChannelId primitive.ObjectID `json:"_id" bson:"_id"`
}

type UploadMediaResponse struct {
	FileName  string `bson:"file_name" json:"file_name" structs:"file_name"`
	SignedUrl string `bson:"signed_url" json:"signed_url" structs:"signed_url"`
}

type GetMessagesByChannelIdResponse struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" structs:"_id"`
	Content     string             `json:"content,omitempty" bson:"content,omitempty" structs:"content"`
	ContentType string             `json:"content_type" bson:"content_type" validate:"required,oneof=text media both" structs:"content_type"`
	// MediaUrl    string             `json:"media_url,omitempty" bson:"media_url,omitempty" structs:"media_url"`
	// ParentId  primitive.ObjectID `json:"parent_id,omitempty" bson:"parent_id,omitempty" structs:"parent_id"`
	UpdatedAt time.Time    `json:"updated_at" bson:"updated_at" structs:"updated_at"`
	User      UserResponse `json:"user" bson:"user" structs:"user"`
}

type MessageResponse struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id" structs:"_id"`
	Content     string             `json:"content" bson:"content" structs:"content"`
	ContentType string             `json:"content_type" bson:"content_type" structs:"content_type"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at" structs:"updated_at"`
	User        AllUserResponse    `json:"user" bson:"user" structs:"user"`
}

type MetaPagination struct {
	Page      int32 `json:"page" bson:"page" structs:"page"`
	TotalPage int32 `json:"total_page" bson:"total_page" structs:"total_page"`
}

func MessageWithCode(status int, message string) map[string]interface{} {
	return map[string]interface{}{"res_code": status, "message": message}
}

func ResponseErrorWithCode(status int, message string) map[string]interface{} {
	return MessageWithCode(status, message)
}

func GetHTTPStatusCode(resCode interface{}) int {
	if resCode != nil {
		return resCode.(int)
	}
	return http.StatusOK
}

func ResponseSuccessWithToken(message string, code int, resData map[string]interface{}) map[string]interface{} {
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: message,
			Code:    code,
		},
	}
	if token, ok := resData["token"]; ok {
		response.Meta.Token = token.(string)
	}
	if rData, ok := resData["data"]; ok {
		response.Data = rData

	} else {
		response.Data = nil
	}

	m := structs.Map(response)

	return m
}

func ConvertToInterface(message string, code int, data interface{}) map[string]interface{} {
	d := map[string]interface{}{
		"message": message,
		"code":    code,
		"data":    data,
	}
	d = FinalResponse(d)
	return d
}

func Response(resData interface{}) map[string]interface{} {
	data := resData.(map[string]interface{})
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: data["message"].(string),
			Code:    data["code"].(int),
		},
	}
	if resData != nil {
		if rData, ok := data["data"]; ok {
			response.Data = rData

		} else {
			response.Data = nil
		}
	} else {
		response.Data = nil
	}

	m := structs.Map(response)
	return m
}

func ResponseSuccessWithCode(message string, data ...interface{}) map[string]interface{} {
	return ConvertToInterface(message, META_SUCCESS, data)
}

func FinalResponse(data map[string]interface{}) map[string]interface{} {
	response := BaseSuccessResponse{
		Meta: Meta{
			Message: data["message"].(string),
			Code:    data["code"].(int),
		},
	}

	if rData, ok := data["data"]; ok {
		response.Data = rData

	} else {
		response.Data = nil
	}

	m := structs.Map(response)
	return m
}

func Response_SignIn(c *gin.Context, status int, data map[string]interface{}) {
	if status != 200 {
		data = FinalResponse(data)
	}
	c.JSON(status, data)
}

func Respond(c *gin.Context, status int, data map[string]interface{}) {
	d := FinalResponse(data)
	c.JSON(status, d)
}
