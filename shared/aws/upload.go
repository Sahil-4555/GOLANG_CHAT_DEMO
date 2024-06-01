package aws

import (
	"time"

	"chat-demo-golang/configs"
	"chat-demo-golang/shared/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(key string) (string, error) {
	accessKey := configs.AccessKey()
	secret := configs.SecretKey()
	region := configs.Region()
	bucket := configs.Bucket()

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secret, ""),
		Region:      aws.String(region),
	})

	if err != nil {
		log.GetLog().Error("ERROR: ", err.Error())
		return "", err
	}

	svc := s3.New(sess)
	resp, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	url, err := resp.Presign(15 * time.Minute)
	if err != nil {
		log.GetLog().Error("ERROR: ", err.Error())
		return "", err
	}
	return url, nil
}
