package service

import (
	"strings"
	"time"

	"github.com/Sahil-4555/mvc/shared/aws"
	"github.com/Sahil-4555/mvc/shared/common"
	"github.com/Sahil-4555/mvc/shared/log"
	"github.com/Sahil-4555/mvc/shared/message"
)

func UploadMedia(req common.UploadMediaReq) map[string]interface{} {

	var respData common.UploadMediaResponse

	if req.FileData != nil {
		err := aws.UploadToS3(req.Key, req.FileSize, req.FileData)
		if err != nil {
			log.GetLog().Info("ERROR : ", err.Error())
			var data interface{}
			return common.ConvertToInterface(message.SomethingWrong, common.META_FAILED, data)
		}

	}
	fileName := strings.Split(req.Key, "/")
	respData.FileName = fileName[len(fileName)-1]

	url, err := aws.GenerateSignedUrl(req.Key, time.Duration(3))
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		var data interface{}
		return common.ConvertToInterface(err.Error(), common.META_FAILED, data)
	}

	respData.SignedUrl = url

	response := common.ResponseSuccessWithCode(message.FileUploadSuccess, respData)
	return response
}
