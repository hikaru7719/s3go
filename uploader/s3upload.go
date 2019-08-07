package uploader

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hikaru7719/s3go/time"
)

var (
	baseHost = "s3.amazonaws.com"
)

// New returns S3Upload
func New(bucketName, fileName string, signature Signature) (*S3Upload, error) {
	host := fmt.Sprintf("%s.%s", bucketName, baseHost)
	return &S3Upload{
		host:       host,
		bucketName: bucketName,
		objectName: fileName,
		signature:  signature,
	}, nil
}

// Signature is interface
type Signature interface {
	Authorization(method, URL, payload string, header map[string]string) string
}

// S3Upload is struct for upliading file to AWS S3
type S3Upload struct {
	host       string
	bucketName string
	objectName string
	uploadID   string
	signature  Signature
	file       io.Reader
}

// InitialMultipartUpload is first request to do maltipart upload
func (s *S3Upload) InitialMultipartUpload() error {
	client := &http.Client{}
	req, err := s.newInitialRequest()
	res, err := client.Do(req)
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	s.xmlMapping(byteBody)
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

type initialRespXML struct {
	UploadID string `xml:"UploadId"`
}

func (s *S3Upload) xmlMapping(respBody []byte) {
	xmlMapper := initialRespXML{}
	xml.Unmarshal(respBody, &xmlMapper)
	s.uploadID = xmlMapper.UploadID
}

func (s *S3Upload) convertToMap(header http.Header) map[string]string {
	newMap := make(map[string]string)
	for key := range header {
		newMap[key] = header.Get(key)
	}
	return newMap
}

// PutMaltiPartObject is request to upload object
func (s *S3Upload) PutMaltiPartObject() error {
	client := &http.Client{}
	req, err := s.newUploaderRequest()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	byteBody, err := ioutil.ReadAll(res.Body)
	fmt.Println(res.StatusCode, string(byteBody))
	return err
}

func (s *S3Upload) newUploaderRequest() (*http.Request, error) {
	url := fmt.Sprintf("https://%s/%s?partNumber=1&uploadId=%s", s.host, s.objectName, s.uploadID)
	byteBody, _ := ioutil.ReadAll(s.file)
	buffer := bytes.NewBuffer(byteBody)
	req, err := http.NewRequest("PUT", url, buffer)
	req.Header.Add("x-amz-date", time.Default.Now())
	req.Header.Add("Host", s.host)
	req.Header.Add("x-amz-content-sha256", hashSHA256(string(byteBody)))
	req.Header.Add("Content-Length", strconv.Itoa(len(byteBody)))
	headerMap := s.convertToMap(req.Header)
	authorization := s.signature.Authorization("PUT", url, string(byteBody), headerMap)
	req.Header.Add("Authorization", authorization)
	return req, err
}

func hashSHA256(payload string) string {
	hash := sha256.Sum256([]byte(payload))
	hexed := hex.EncodeToString(hash[:])
	return strings.ToLower(hexed)
}
