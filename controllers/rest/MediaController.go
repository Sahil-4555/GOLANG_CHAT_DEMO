package rest

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/Sahil-4555/mvc/configs/middleware"
	service "github.com/Sahil-4555/mvc/services"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
	"github.com/Sahil-4555/mvc/shared/utils"
	validators "github.com/Sahil-4555/mvc/validators"
	"github.com/gin-gonic/gin"
)

// UploadMedia
// @Summary UploadMedia - This API is used to upload media files.
// @tags MediaMethod
// @security BearerAuth
// @Param channelId path string true "Channel ID"
// @Param media formData file true "Media file to upload"
// @router /v1/media/uploadmedia/:channelId [post]
func UploadMedia(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Media Controller Called(UploadMedia).")

	var req common.UploadMediaReq
	channelId := c.Param("channelId")
	userId, err := middleware.GetUserData(c)
	if err != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(err.Error(), common.META_FAILED, data))
		return
	}

	maxFileSize := 21474836478989892
	req.UserId = userId

	mediaFile, handler, err := c.Request.FormFile("media")

	if err != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(err.Error(), common.META_FAILED, data))
		return
	}

	basePath := channelId

	req.FileSize = handler.Size
	data, err := createMediaRequest(handler, mediaFile, int64(maxFileSize), utils.GenerateID(), basePath)

	if err == nil {
		req.Key = data["key"].(string)
		req.FileData = data["data"].([]byte)
		req.FileName = data["name"].(string)
	} else {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.GeneralInvalidFile, common.META_FAILED, data))
		return
	}

	defer func(mediaFile multipart.File) {
		err := mediaFile.Close()
		if err != nil {
			log.GetLog().Error(err, "ERROR(Error while closing file) : ")
		}
	}(mediaFile)

	if resp, ok := validators.ValidateStruct(req, "UploadMediaReq"); !ok {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(resp, common.META_FAILED, data))
		return
	}

	resp := service.UploadMedia(req)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Media Uploded successfully...")
}

func createMediaRequest(header *multipart.FileHeader, fileImg multipart.File, maxSize int64, uuId string, basePath string) (map[string]interface{}, error) {
	if header.Size <= maxSize {
		extension := header.Filename
		extArr := strings.Split(extension, ".")
		fileName := uuId + "." + extArr[1]
		key := basePath + "/" + fileName
		data, _ := io.ReadAll(fileImg)
		response := map[string]interface{}{
			"key":  key,
			"data": data,
			"name": fileName,
		}
		return response, nil
	} else {
		return map[string]interface{}{}, errors.New(message.InvalidSizeOrType)
	}
}
