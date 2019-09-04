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
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/hikaru7719/s3go/time"
)

var (
	baseHost = "s3.amazonaws.com"
)

// New returns S3Upload
func New(bucketName, fileName string, signature Signature) (*S3Upload, error) {
	host := fmt.Sprintf("%s.%s", bucketName, baseHost)
	etagMapper := make(map[int]string, 20)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	mutex := new(sync.Mutex)
	return &S3Upload{
		host:       host,
		bucketName: bucketName,
		objectName: fileName,
		signature:  signature,
		etagMapper: etagMapper,
		file:       file,
		mutex:      mutex,
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
	file       io.ReadCloser
	etagMapper map[int]string
	fileSlice  [][]byte
	mutex      *sync.Mutex
}

// Run runs to upload file
func (s *S3Upload) Run() error {
	err := s.InitialMultipartUpload()
	if err != nil {
		return err
	}
	err = s.devideFile()
	if err != nil {
		return err
	}
	s.PutObject()
	err = s.CompleteUploadObject()
	if err != nil {
		return err
	}
	defer s.file.Close()
	return nil
}

// InitialMultipartUpload is first request to do maltipart upload
func (s *S3Upload) InitialMultipartUpload() error {
	client := &http.Client{}
	req, err := s.newInitialRequest()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	byteBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	s.xmlMapping(byteBody)
	return nil
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

// PutObject uploads file divided some chunk
func (s *S3Upload) PutObject() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	errChan := make(chan error)

	go func() {
		for err := range errChan {
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	for n := range s.fileSlice {
		wg.Add(1)
		go func(n int) {
			s.PutMultiPartObject(n+1, done, errChan)
			defer wg.Done()
		}(n)
	}

	wg.Wait()
	close(errChan)
	close(done)
}

// PutMultiPartObject is request to upload object
func (s *S3Upload) PutMultiPartObject(partNumber int, done <-chan interface{}, errChan chan<- error) {
	fmt.Println("start upload part", partNumber)
	select {
	case <-done:
		return
	default:
	}

	client := &http.Client{}
	req, err := s.newUploaderRequest(partNumber)
	res, err := client.Do(req)
	if err != nil {
		errChan <- err
	}
	defer res.Body.Close()
	etag := res.Header.Get("ETag")
	s.mutexMapInsert(partNumber, etag)
}

func (s *S3Upload) mutexMapInsert(partNumber int, etag string) {
	s.mutex.Lock()
	s.etagMapper[partNumber] = etag
	defer s.mutex.Unlock()
}

func (s *S3Upload) devideFile() error {
	byteSlice := make([][]byte, 0, 10)
	for {
		bytes := make([]byte, 1024*1024*5)
		size, err := s.file.Read(bytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err

		}
		if size != 1024*1024*5 {
			byteSlice = append(byteSlice, bytes[:size])
			continue
		}
		byteSlice = append(byteSlice, bytes)
	}
	s.fileSlice = byteSlice
	return nil
}

func (s *S3Upload) newUploaderRequest(partNumber int) (*http.Request, error) {
	url := fmt.Sprintf("https://%s/%s?partNumber=%d&uploadId=%s", s.host, s.objectName, partNumber, s.uploadID)
	byteBody := s.fileSlice[partNumber-1]
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

func (s *S3Upload) etagMapping(partNumber int, etag string) {
	s.etagMapper[partNumber] = etag
}

func hashSHA256(payload string) string {
	hash := sha256.Sum256([]byte(payload))
	hexed := hex.EncodeToString(hash[:])
	return strings.ToLower(hexed)
}

// CompleteUploadObject is request to finish upload part
func (s *S3Upload) CompleteUploadObject() error {
	client := &http.Client{}
	req, err := s.newCompleteRequest()
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	byteBody, err := ioutil.ReadAll(res.Body)
	fmt.Println(res.StatusCode, string(byteBody))
	return err
}

// CompleteMultipartUpload struct is to be base XML For Reuqest
type CompleteMultipartUpload struct {
	XMLName xml.Name `xml:"CompleteMultipartUpload"`
	Part    []Part
}

// Part is included in the CompleteMultipartUpload XML.
type Part struct {
	XMLName    xml.Name `xml:"Part"`
	PartNumber int      `xml:"PartNumber"`
	ETag       string   `xml:"ETag"`
}

// Parts implements Sort Interface
type Parts []Part

func (p Parts) Len() int           { return len(p) }
func (p Parts) Less(i, j int) bool { return p[i].PartNumber < p[j].PartNumber }
func (p Parts) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (s *S3Upload) generateXML() (string, error) {
	parts := make([]Part, 0, 10)
	for key, value := range s.etagMapper {
		part := Part{PartNumber: key, ETag: value}
		parts = append(parts, part)
	}
	sort.Sort(Parts(parts))
	comleteMultipartUpload := CompleteMultipartUpload{Part: parts}
	bytes, err := xml.MarshalIndent(comleteMultipartUpload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *S3Upload) newCompleteRequest() (*http.Request, error) {
	url := fmt.Sprintf("https://%s/%s?uploadId=%s", s.host, s.objectName, s.uploadID)
	xmlString, err := s.generateXML()
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(xmlString)
	req, err := http.NewRequest("POST", url, reader)
	req.Header.Add("x-amz-date", time.Default.Now())
	req.Header.Add("Host", s.host)
	req.Header.Add("x-amz-content-sha256", hashSHA256(xmlString))
	req.Header.Add("Content-Length", strconv.Itoa(len(xmlString)))
	headerMap := s.convertToMap(req.Header)
	authorization := s.signature.Authorization("POST", url, xmlString, headerMap)
	req.Header.Add("Authorization", authorization)
	return req, err
}
