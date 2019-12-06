package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
	"log"
	"os"
)

//NewSession creates a new session for file upload
func NewSession() *session.Session {
	sess, err := session.NewSession(newConfig())
	if err != nil {
		log.Printf("Unable to create session: %+v", err)
	}
	return sess
}

//
func newConfig() *aws.Config {
	return &aws.Config{
		Region:      aws.String(viper.GetString("region")),
		Credentials: credentials.NewStaticCredentials(viper.GetString("aws_access_key"), viper.GetString("aws_secret_key"), ""),
	}
}

//Upload file to s3
func Upload(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Unable to read file: %+v", err)
		return
	}
	uploader := s3manager.NewUploader(NewSession())
	key := filePath
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(viper.GetString("bucket")),
		Key:    &key,
		Body:   file,
	})

	if err != nil {
		log.Printf("Failed to upload file, %v", err)
		return
	}

	log.Printf("Uploaded file to S3: %+v", result)

	if err := file.Close(); err != nil {
		log.Printf("Unable to close file reader: %+v", err)
		return
	}

	if err != nil {
		log.Printf("Failed to upload file, %v", err)
	} else {
		cleanUp(file.Name())
	}
}

func cleanUp(dir string) {
	if err := os.Remove(dir); err != nil {
		log.Printf("Unable to clean up file: %+v", err)
	} else {
		log.Printf("Deleted: %s", dir)
	}
}
