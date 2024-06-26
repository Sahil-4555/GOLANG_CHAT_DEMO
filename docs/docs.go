// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/api/login": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "UserMethod"
                ],
                "summary": "SignIn - This API is used to SignIn the user.",
                "parameters": [
                    {
                        "description": "signin request",
                        "name": "SignupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.SignInReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/api/signup": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "UserMethod"
                ],
                "summary": "SignUp - This API is used to SignUp the user.",
                "parameters": [
                    {
                        "description": "signup request",
                        "name": "SignupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.SignUpReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/addfavoritechannel": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "AddRemoveFavouriteChannel - This API is used to add or remove favourite channel",
                "parameters": [
                    {
                        "description": "add or remove favourite channel",
                        "name": "AddFavoriteChannelRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.AddFavoriteChannelReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/addmembers": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "AddMembersToGroup - This API is used to add members to group channel",
                "parameters": [
                    {
                        "description": "add members to group channel",
                        "name": "AddMembersToGroupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.AddMembersToGroupReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/closeconversation": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "CloseConversation - This API is used to close conversation for one to one channels",
                "parameters": [
                    {
                        "description": "close conversation for one to one channels",
                        "name": "CloseConversationRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.CloseConversationReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/createGroup": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "CreateGroup - This API is used to create a new group",
                "parameters": [
                    {
                        "description": "create group",
                        "name": "CreateGroupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.CreateGroupReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/get-all-direct-channel": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GetAllOneToOneChannelInOrder - This API is used to get all one to one channels of user in order",
                "responses": {}
            }
        },
        "/v1/channel/get-direct-channel": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GetOneToOneChannelsConnectedWithUser - This API is used to get all one to one channels connected with user",
                "responses": {}
            }
        },
        "/v1/channel/get-favourite-channel": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GetFavouritesChannel - This API is used to get all favourite channels",
                "responses": {}
            }
        },
        "/v1/channel/get-group-channel": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GetPrivateChannelsConnectedWithUser - This API is used to get all private channels connected with user",
                "responses": {}
            }
        },
        "/v1/channel/getrecentchannels": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GetRecentChannels - This API is used to get recent channels of users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Offset value",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/give-admin-rights": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "GiveAdminRightsToUser - This API is used to give admin rights to user",
                "parameters": [
                    {
                        "description": "give admin rights to user",
                        "name": "GiveAdminRightsToUserRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.GiveAdminRightsToUserReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/join-group/{channelId}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "JoinGroup - This API is used to join group channel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Channel ID",
                        "name": "channelId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/joinRoom": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "JoinChannel - This API is used to join a channel",
                "parameters": [
                    {
                        "description": "join channel",
                        "name": "JoinChannelRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.JoinChannelReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/leavechannel": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "LeaveChannel - This API is used to leave group channel",
                "parameters": [
                    {
                        "description": "leave group channel",
                        "name": "LeaveChannelRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.LeaveChannelReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/channel/removeuserbyadmin": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChannelMethod"
                ],
                "summary": "AddMembersToGroup - This API is used to remove members from group channel by group admin",
                "parameters": [
                    {
                        "description": "remove members from group channel by group admin",
                        "name": "RemoveMemberFromGroupRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.RemovUserFromGroupByGroupAdminReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/deletemessage": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "DeleteMessage - This API is used to delete the chat message.",
                "parameters": [
                    {
                        "description": "delete message",
                        "name": "DeleteMessageRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.DeleteMessageReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/get-all-users": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "GetAllUsers - This API is used to search the users based on search value",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search Value",
                        "name": "searchValue",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/get-channel-members": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "GetChannelMembers - This API is used to get the channel members of channel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Channel Id",
                        "name": "channelId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/messages-by-channelid/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "GetMessagesById - This API is used to get the messages by ChannelId",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Channel Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Page Value",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Offset Value",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/searchhandler": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "SearchHandler - This API is used to search channels and users based on search value.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search Value",
                        "name": "searchValue",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/v1/chat/updatemessage": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "ChatMethod"
                ],
                "summary": "UpdateMessage - This API is used to update the chat message.",
                "parameters": [
                    {
                        "description": "update message",
                        "name": "UpdateMessageRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/common.UpdateMessageReq"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/media/uploadmedia/:channelId": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "MediaMethod"
                ],
                "summary": "UploadMedia - This API is used to upload media files.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Channel ID",
                        "name": "channelId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Media file to upload",
                        "name": "media",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "common.AddFavoriteChannelReq": {
            "type": "object",
            "properties": {
                "channel_id": {
                    "type": "string"
                },
                "is_favourite": {
                    "type": "boolean"
                }
            }
        },
        "common.AddMembersToGroupReq": {
            "type": "object",
            "required": [
                "channel_id",
                "users"
            ],
            "properties": {
                "channel_id": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "common.CloseConversationReq": {
            "type": "object",
            "required": [
                "channel_id"
            ],
            "properties": {
                "channel_id": {
                    "type": "string"
                }
            }
        },
        "common.CreateGroupReq": {
            "type": "object",
            "required": [
                "channel_name"
            ],
            "properties": {
                "channel_name": {
                    "type": "string",
                    "maxLength": 50
                },
                "description": {
                    "type": "string"
                }
            }
        },
        "common.DeleteMessageReq": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                }
            }
        },
        "common.GiveAdminRightsToUserReq": {
            "type": "object",
            "properties": {
                "channel_id": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "common.JoinChannelReq": {
            "type": "object",
            "properties": {
                "reciver_id": {
                    "type": "string"
                }
            }
        },
        "common.LeaveChannelReq": {
            "type": "object",
            "required": [
                "channel_id"
            ],
            "properties": {
                "channel_id": {
                    "type": "string"
                }
            }
        },
        "common.RemovUserFromGroupByGroupAdminReq": {
            "type": "object",
            "required": [
                "channel_id",
                "users"
            ],
            "properties": {
                "channel_id": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "common.SignInReq": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "common.SignUpReq": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "user_name"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string",
                    "maxLength": 50
                },
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string",
                    "maxLength": 50
                }
            }
        },
        "common.UpdateMessageReq": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string"
                },
                "content": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Bearer token authentication",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Chat Application API",
	Description:      "Chat Application API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
