package uploader

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hikaru7719/s3go/time"
)

var (
	BASE_HOST = "s3.amazonaws.com"
)

func New(bucketName, fileName string, signature Signature) (*S3Upload, error) {
	host := fmt.Sprintf("%s.%s", bucketName, BASE_HOST)
	return &S3Upload{
		host:       host,
		bucketName: bucketName,
		objectName: fileName,
		signature:  signature,
	}, nil
}

type Signature interface {
	Authorization(method, URL, payload string, header map[string]string) string
}

type S3Upload struct {
	host       string
	bucketName string
	objectName string
	uploadId   string
	signature  Signature
}

func (s *S3Upload) InitialMultipartUpload() error {
	client := &http.Client{}
	req, err := s.newInitialRequest()
	res, err := client.Do(req)
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	fmt.Println(res.StatusCode, string(byteBody))
	return err
}

func (s *S3Upload) newInitialRequest() (*http.Request, error) {
	url := fmt.Sprintf("https://%s/%s?uploads", s.host, s.objectName)
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("x-amz-date", time.Default.Now())
	req.Header.Add("Host", s.host)
	req.Header.Add("x-amz-content-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	headerMap := s.convertToMap(req.Header)
	authorization := s.signature.Authorization("POST", url, "", headerMap)
	req.Header.Add("Authorization", authorization)
	return req, err
}

func (s *S3Upload) convertToMap(header http.Header) map[string]string {
	newMap := make(map[string]string)
	for key := range header {
		newMap[key] = header.Get(key)
	}
	return newMap
}
