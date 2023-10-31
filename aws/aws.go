package aws

import (
	"bytes"
	"fmt"

	"github.com/Sync-Space-49/syncspace-server/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3Client *s3.S3
)

func GetS3Client() error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	if s3Client == nil {
		// https://knowledgebase.wasabi.com/hc/en-us/articles/360000762391-How-do-I-use-AWS-SDK-for-Go-Golang-with-Wasabi-
		s3Config := aws.Config{
			Credentials:      credentials.NewStaticCredentials(cfg.Wasabi.AccessKey, cfg.Wasabi.SecretKey, ""),
			Endpoint:         aws.String("https://s3.wasabisys.com"),
			Region:           aws.String("us-east-1"),
			S3ForcePathStyle: aws.Bool(true),
		}

		goSession, err := session.NewSessionWithOptions(session.Options{
			Config: s3Config,
		})

		if err != nil {
			return err
		}
		s3Client = s3.New(goSession)
	}

	return nil
}

func UploadPfp(file *bytes.Reader, filename string) (*string, error) {
	cfg, err := config.Get()
	if err != nil {
		return nil, err
	}

	if s3Client == nil {
		err := GetS3Client()
		if err != nil {
			return nil, err
		}
	}

	filepath := fmt.Sprintf("pfp/%s", filename)
	putObjectInput := &s3.PutObjectInput{
		Body:        file,
		Bucket:      aws.String(cfg.Wasabi.Bucket),
		Key:         aws.String(filepath),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/png"),
	}

	_, err = s3Client.PutObject(putObjectInput)
	if err != nil {
		return nil, err
	}

	fileurl := fmt.Sprintf("https://s3.%s.wasabisys.com/%s/%s", cfg.Wasabi.Region, cfg.Wasabi.Bucket, filepath)
	return &fileurl, nil
}
