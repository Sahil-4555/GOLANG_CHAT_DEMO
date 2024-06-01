package aws

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/Sahil-4555/mvc/configs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadToS3(key string, size int64, buffer []byte) error {
	accessKey := configs.AccessKey()
	secret := configs.SecretKey()
	region := configs.Region()
	bucket := configs.Bucket()
	s, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secret, ""),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		return err
	}

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("inline"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}
	return nil
}

func GenerateSignedUrl(keyPath string, minutes time.Duration) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(configs.Region()),
	})
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(configs.Bucket()),
		Key:    aws.String(keyPath),
	})
	urlStr, err := req.Presign(minutes * time.Hour)

	if err != nil {
		fmt.Println("Failed to sign request", err)
		return "", err
	}

	return urlStr, nil
}
