package uploader

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hikaru7719/s3go/config"
	"github.com/hikaru7719/s3go/signature"
	"github.com/hikaru7719/s3go/time"
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
	uploadId   string
}

func (s *S3Upload) InitialMultipartUpload() error {
	client := &http.Client{}
	host := fmt.Sprintf("%s.%s", s.bucketName, BASE_HOST)
	url := fmt.Sprintf("https://%s/%s?uploads", host, s.objectName)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("x-amz-date", time.Default.Now())
	req.Header.Add("Host", host)
	req.Header.Add("x-amz-content-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	headerMap := s.convertToMap(req.Header)
	authorization := signature.Authorization("POST", url, "", headerMap, config.Default, time.Default)
	req.Header.Add("Authorization", authorization)
	res, err := client.Do(req)
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(res.StatusCode, string(byteBody))
	return err
}

func (s *S3Upload) convertToMap(header http.Header) map[string]string {
	newMap := make(map[string]string)
	for key := range header {
		newMap[key] = header.Get(key)
	}
	return newMap
}
