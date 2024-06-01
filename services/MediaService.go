package service

import (
	"errors"
	"strings"

	"chat-demo-golang/shared/aws"
	"chat-demo-golang/shared/common"
	"chat-demo-golang/shared/log"
	"chat-demo-golang/shared/utils"
)

const (
	maxFileSize = 10 * 1024 * 1024
)

func CreateMediaRequest(filesize int64, filename string, uuId string, basePath string) (map[string]interface{}, error) {
	log.GetLog().Info("INFO : ", "S3 Service Called(createMediaRequest).")
	if filesize <= maxFileSize {
		extension := filename
		extArr := strings.Split(extension, ".")
		fileName := extArr[0] + "-" + utils.GenerateID() + "." + extArr[1]
		key := basePath + "/" + fileName
		response := map[string]interface{}{
			"key":  key,
			"name": fileName,
		}
		return response, nil
	} else {
		return map[string]interface{}{}, errors.New("Invalid file type or size")
	}
}

func UploadMedia(req common.UploadMediaData) map[string]interface{} {
	log.GetLog().Info("INFO : ", "S3 Service Called(UploadMedia).")
	var respData common.UploadMediaResponse

	url, err := aws.UploadToS3(req.Key)
	if err != nil {
		log.GetLog().Error("ERROR : ", err.Error())
		return map[string]interface{}{
			"message":  err.Error(),
			"code":     common.META_FAILED,
			"res_code": common.STATUS_OK,
		}
	}

	fileName := strings.Split(req.Key, "/")
	respData.FileName = fileName[len(fileName)-1]
	respData.PreSignedUrl = url

	return map[string]interface{}{
		"message": "File uploaded successfully.",
		"code":    common.META_SUCCESS,
		"data":    respData,
	}
}
