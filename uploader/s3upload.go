package uploader

import (
	"fmt"
	"net/http"
)

var (
	BASE_HOST = "s3.amazonaws.com"
)

func New(bucketName, fileName string) (*S3Upload, error) {
	return &S3Upload{
		bucketName: bucketName,
		objectName: fileName,
	}, nil
}

type S3Upload struct {
	bucketName string
	objectName string
}

func (s *S3Upload) RequestS3() error {
	client := &http.Client{}
	host := fmt.Sprintf("https://%s.%s", s.bucketName, BASE_HOST)
	req, err := http.NewRequest("POST", host, nil)
	client.Do(req)
	return err
}
