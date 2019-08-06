package uploader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hikaru7719/s3go/signature"
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
	url        string
}

func (s *S3Upload) InitialMultipartUpload() error {
	client := &http.Client{}
	host := fmt.Sprintf("%s.%s", s.bucketName, BASE_HOST)
	url := fmt.Sprintf("https://%s/%s?uploads", host, s.objectName)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("X-Amz-Date", time.Now().Format("20060102'T'150405'Z'"))
	req.Header.Add("Host", host)
	headerMap := s.convertToMap(req.Header)
	for key, value := range headerMap {
		fmt.Println(key, ":", value)
	}

	authorization := signature.Authorization("POST", url, "", headerMap)
	fmt.Println(authorization)
	req.Header.Add("Authorization", authorization)
	res, err := client.Do(req)
	byteBody, err := ioutil.ReadAll(res.Body)
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
