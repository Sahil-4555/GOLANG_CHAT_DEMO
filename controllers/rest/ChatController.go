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

func SearchHandler(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(SearchHandler).")

	userId, _ := middleware.GetUserData(c)
	searchValue := c.Query("searchValue")

	var resp map[string]interface{}
	if len(searchValue) < 3 {
		resp = service.GetAllOneToOneChannelConnectedWithUserInOrder(userId)
	} else {
		resp = service.SearchHandler(searchValue, userId)
	}

	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Searched successfully...")
}

func UpdateMessage(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(UpdateMessage).")

	var req common.UpdateMessageReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "UpdateMessageReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.UpdateMessage(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Message updated successfully...")
}

func DeleteMessage(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(DeleteMessage).")

	var req common.DeleteMessageReq
	userId, _ := middleware.GetUserData(c)

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_SUCCESS, data))
		return
	}

	if resp, ok := validators.ValidateStruct(req, "DeleteMessageReq"); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": map[string]interface{}{
				"code":    common.META_FAILED,
				"message": resp,
			},
		})
		return
	}

	resp := service.DeleteMessage(req, userId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Message deleted successfully...")
}

func GetMessageById(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(GetMessageById).")

	Id := c.Param("_id")
	messageId, _ := primitive.ObjectIDFromHex(Id)

	resp := service.GetMessageById(messageId)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Message fetched successfully...")
}

func GetMessagesByChannelId(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(GetMessagesByChannelId).")

	Id := c.Param("_id")
	channelId, _ := primitive.ObjectIDFromHex(Id)

	pageStr := c.Query("page")
	page, _ := strconv.Atoi(pageStr)

	offsetStr := c.Query("offset")
	offset, _ := strconv.Atoi(offsetStr)

	resp := service.GetMessagesByChannelId(channelId, page, offset)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Messages fetched successfully...")
}

func GetAllUsers(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(GetAllUsers).")

	searchValue := c.Query("searchValue")

	var resp map[string]interface{}
	if len(searchValue) < 3 {
		resp = service.GetAllUsers()
	} else {
		resp = service.SearchUser(searchValue)
	}
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Users fetched successfully...")
}

func GetChannelMembers(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Chat Controller Called(GetChannelMembers).")

	channelId := c.Query("channelId")
	Id, _ := primitive.ObjectIDFromHex(channelId)

	resp := service.GetChannelMembers(Id)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Channel members fetched successfully...")
}
