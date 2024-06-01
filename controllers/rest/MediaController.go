package rest

import (
	"net/http"

	service "chat-demo-golang/services"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"chat-demo-golang/shared/message"
	"chat-demo-golang/shared/utils"

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
	log.GetLog().Info("INFO : ", "S3 Controller Called(UploadMedia).")

	var req common.UploadMediaRequest

	if c.BindJSON(&req) != nil {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.FailedToReadBody, common.META_FAILED, data))
		return
	}

	basePath := req.ChannelId
	data, err := service.CreateMediaRequest(req.FileSize, req.FileName, utils.GenerateID(), basePath)

	var filedata common.UploadMediaData
	if err == nil {
		filedata.Key = data["key"].(string)
		filedata.FileName = data["name"].(string)
	} else {
		var data interface{}
		common.Respond(c, http.StatusBadRequest, common.ConvertToInterface(message.GeneralInvalidFile, common.META_FAILED, data))
		return
	}

	resp := service.UploadMedia(filedata)
	statusCode := common.GetHTTPStatusCode(resp["res_code"])
	common.Respond(c, statusCode, resp)
	log.GetLog().Info("INFO : ", "Media Uploded successfully...")
}

/*
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:PutObject",
                "s3:GetObject"
            ],
            "Resource": "arn:aws:s3:::chat-demo-golang/*"
        }
    ]
}

[
    {
        "AllowedHeaders": [
            "*"
        ],
        "AllowedMethods": [
            "GET",
            "PUT",
            "POST",
            "DELETE"
        ],
        "AllowedOrigins": [
            "*"
        ],
        "ExposeHeaders": [
            "ETag"
        ],
        "MaxAgeSeconds": 3000
    }
]
*/
