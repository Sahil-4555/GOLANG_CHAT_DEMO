package rest

import (
	"net/http"
	"strconv"

	"github.com/Sahil-4555/mvc/configs/middleware"
	service "github.com/Sahil-4555/mvc/services"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	validators "github.com/Sahil-4555/mvc/validators"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

func GetAllOneToOneChannelConnectedWithUserInOrder(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Channel Controller Called(GetAllOneToOneChannelConnectedWithUserInOrder).")

	userId, _ := middleware.GetUserData(c)

	resp := service.GetAllOneToOneChannelConnectedWithUserInOrder(userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "One to one channel in order fetched successfully.")
}

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
