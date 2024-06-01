package rest

import (
	"net/http"
	"strconv"

	"chat-demo-golang/configs/middleware"
	service "chat-demo-golang/services"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"chat-demo-golang/shared/message"
	validators "chat-demo-golang/validators"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JoinChannel
// @Summary JoinChannel - This API is used to join a channel
// @tags ChannelMethod
// @security BearerAuth
// @Param JoinChannelRequest body common.JoinChannelReq true "join channel"
// @router /v1/channel/joinRoom [post]
func JoinChannel(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(JoinChannel).")

	var req common.JoinChannelReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "JoinChannelReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.JoinChannel(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Join channel successfully...")
}

// GetOneToOneChannelsConnectedWithUser
// @Summary GetOneToOneChannelsConnectedWithUser - This API is used to get all one to one channels connected with user
// @tags ChannelMethod
// @security BearerAuth
// @router /v1/channel/get-direct-channel [get]
func GetOneToOneChannelsConnectedWithUser(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetOneToOneChannelsConnectedWithUser).")

	userId, err := middleware.GetUserData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": err.Error(),
			},
		})
		return
	}

	resp := service.GetOneToOneChannelsConnectedWithUser(userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "One to one channel fetched successfully...")
}

// GetPrivateChannelsConnectedWithUser
// @Summary GetPrivateChannelsConnectedWithUser - This API is used to get all private channels connected with user
// @tags ChannelMethod
// @security BearerAuth
// @router /v1/channel/get-group-channel [get]
func GetPrivateChannelsConnectedWithUser(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetPrivateChannelsConnectedWithUser).")

	userId, err := middleware.GetUserData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": err.Error(),
			},
		})
		return
	}

	resp := service.GetPrivateChannelsConnectedWithUser(userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Group Channel fetched successfully...")
}

// CreateGroup
// @Summary CreateGroup - This API is used to create a new group
// @tags ChannelMethod
// @security BearerAuth
// @Param CreateGroupRequest body common.CreateGroupReq true "create group"
// @router /v1/channel/createGroup [post]
func CreateGroup(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(CreateGroup).")

	var req common.CreateGroupReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "CreateGroupReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.CreateGroup(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Created group successfully.")
}

// AddRemoveFavouriteChannel
// @Summary AddRemoveFavouriteChannel - This API is used to add or remove favourite channel
// @tags ChannelMethod
// @security BearerAuth
// @Param AddFavoriteChannelRequest body common.AddFavoriteChannelReq true "add or remove favourite channel"
// @router /v1/channel/addfavoritechannel [post]
func AddRemoveFavouriteChannel(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(AddRemoveFavouriteChannel).")

	var req common.AddFavoriteChannelReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "AddRemoveFavouriteChannelReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.AddRemoveFavouriteChannel(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Add/Remove Favourite Channel successfully.")
}

// GetFavouritesChannel
// @Summary GetFavouritesChannel - This API is used to get all favourite channels
// @tags ChannelMethod
// @security BearerAuth
// @router /v1/channel/get-favourite-channel [get]
func GetFavouritesChannel(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetFavouritesChannel).")

	userId, err := middleware.GetUserData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": err.Error(),
			},
		})
		return
	}

	resp := service.GetFavouritesChannel(userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Favourite channels fetched successfully.")
}

// AddMembersToGroup
// @Summary AddMembersToGroup - This API is used to add members to group channel
// @tags ChannelMethod
// @security BearerAuth
// @Param AddMembersToGroupRequest body common.AddMembersToGroupReq true "add members to group channel"
// @router /v1/channel/addmembers [put]
func AddMembersToGroup(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(AddMembersToGroup).")

	var req common.AddMembersToGroupReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "AddMembersToGroupReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.AddMembersToGroup(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Members added to group channel successfully.")
}

// LeaveChannel
// @Summary LeaveChannel - This API is used to leave group channel
// @tags ChannelMethod
// @security BearerAuth
// @Param LeaveChannelRequest body common.LeaveChannelReq true "leave group channel"
// @router /v1/channel/leavechannel [put]
func LeaveChannel(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(LeaveChannel).")

	var req common.LeaveChannelReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "LeaveChannelReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.LeaveChannel(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Leaved channel successfully.")
}

// RemoveUserFromGroupByGroupAdmin
// @Summary AddMembersToGroup - This API is used to remove members from group channel by group admin
// @tags ChannelMethod
// @security BearerAuth
// @Param RemoveMemberFromGroupRequest body common.RemovUserFromGroupByGroupAdminReq true "remove members from group channel by group admin"
// @router /v1/channel/removeuserbyadmin [put]
func RemoveUserFromGroupByGroupAdmin(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(RemoveUserFromGroupByGroupAdmin).")

	var req common.RemovUserFromGroupByGroupAdminReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "RemoveUserFromGroupByGroupAdminReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.RemoveUserFromGroupByGroupAdmin(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Remove user from group channel by admin successfully.")
}

// CloseConversation
// @Summary CloseConversation - This API is used to close conversation for one to one channels
// @tags ChannelMethod
// @security BearerAuth
// @Param CloseConversationRequest body common.CloseConversationReq true "close conversation for one to one channels"
// @router /v1/channel/closeconversation [put]
func CloseConversation(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(CloseConversation).")

	var req common.CloseConversationReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "CloseConversationReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.CloseConversation(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Closed conversation successfully.")
}

// GetRecentChannels
// @Summary GetRecentChannels - This API is used to get recent channels of users
// @tags ChannelMethod
// @security BearerAuth
// @Param page query int false "Page number"
// @Param offset query int false "Offset value"
// @router /v1/channel/getrecentchannels [get]
func GetRecentChannelsOfUsers(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetRecentChannelsOfUsers).")

	userId, _ := middleware.GetUserData(c)

	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)

	offsetStr := c.Query("offset")
	offset, _ := strconv.Atoi(offsetStr)

	resp := service.GetRecentChannelsOfUsers(userId, int64(page), int64(offset))
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Recent channel fetched successfully.")
}

// GetAllOneToOneChannelInOrder
// @Summary GetAllOneToOneChannelInOrder - This API is used to get all one to one channels of user in order
// @tags ChannelMethod
// @security BearerAuth
// @router /v1/channel/get-all-direct-channel [get]
func GetAllOneToOneChannelConnectedWithUserInOrder(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetAllOneToOneChannelConnectedWithUserInOrder).")

	userId, _ := middleware.GetUserData(c)

	resp := service.GetAllOneToOneChannelConnectedWithUserInOrder(userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "One to one channel in order fetched successfully.")
}

// JoinGroup
// @Summary JoinGroup - This API is used to join group channel
// @tags ChannelMethod
// @security BearerAuth
// @Param channelId path string true "Channel ID"
// @router /v1/channel/join-group/{channelId} [get]
func JoinGroupWithChannelId(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(JoinGroupWithChannelId).")

	userId, _ := middleware.GetUserData(c)
	chanId := c.Param("channelId")
	channelId, _ := primitive.ObjectIDFromHex(chanId)

	resp := service.JoinGroupWithChannelId(userId, channelId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Joined group channel by successfully.")
}

// GiveAdminRightsToUser
// @Summary GiveAdminRightsToUser - This API is used to give admin rights to user
// @tags ChannelMethod
// @security BearerAuth
// @Param GiveAdminRightsToUserRequest body common.GiveAdminRightsToUserReq true "give admin rights to user"
// @router /v1/channel/give-admin-rights [put]
func GiveAdminRightsToUser(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GiveAdminRightsToUser).")

	var req common.GiveAdminRightsToUserReq
	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "GiveAdminRightsToUserReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.GiveAdminRightsToUser(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Assigned admin rights to user successfully.")

}
