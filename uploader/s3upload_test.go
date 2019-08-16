package uploader

import (
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/hikaru7719/s3go/signature"
	"github.com/stretchr/testify/assert"
)

func TestInitialMultipartUpload(t *testing.T) {
	sig := signature.New()
	upload, _ := New("s3go-cli-test", "testdata/earth.jpg", sig)
	upload.InitialMultipartUpload()
}

func TestNewInitialRequest(t *testing.T) {

	uploader := &S3Upload{
		host:       "testhost",
		bucketName: "testbucket",
		objectName: "testObject",
		signature:  &mockAuth{},
	}

	req, _ := uploader.newInitialRequest()
	assert.Equal(t, "testhost", req.Header.Get("Host"))
	assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", req.Header.Get("x-amz-content-sha256"))
	assert.Equal(t, "testAuthorization", req.Header.Get("Authorization"))

}

func TestXMLMapping(t *testing.T) {
	uploader := &S3Upload{}
	testXML := `<?xml version="1.0" encoding="UTF-8"?><InitiateMultipartUploadResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Bucket>test</Bucket><Key>test</Key><UploadId>hogehoge</UploadId></InitiateMultipartUploadResult>`
	uploader.xmlMapping([]byte(testXML))
	assert.Equal(t, "hogehoge", uploader.uploadID)
}

type mockAuth struct{}

func (m *mockAuth) Authorization(method, URL, payload string, header map[string]string) string {
	return "testAuthorization"
}

// func TestPutMultipartObject(t *testing.T) {
// 	sig := signature.New()
// 	upload, _ := New("s3go-cli-test", "testdata/earth.jpg", sig)
// 	upload.InitialMultipartUpload()
// 	upload.devideFile()
// 	upload.PutObject()
// 	upload.CompleteUploadObject()
// }

func TestNewUploadRequest(t *testing.T) {
	upload := &S3Upload{
		host:       "testhost",
		bucketName: "testbucket",
		objectName: "testObject",
		signature:  &mockAuth{},
	}
	bytes := []byte("hoge")
	fileSlice := [][]byte{bytes}
	upload.fileSlice = fileSlice
	req, _ := upload.newUploaderRequest(1)
	assert.Equal(t, "4", req.Header.Get("content-length"))
}

func TestNewCompleteRequest(t *testing.T) {
	etagMap := make(map[int]string)
	etagMap[1] = "testetag1"
	upload := &S3Upload{
		host:       "testhost",
		bucketName: "testbucket",
		objectName: "testObject",
		signature:  &mockAuth{},
		etagMapper: etagMap,
	}
	req, _ := upload.newCompleteRequest()
	expectXMLString := `<CompleteMultipartUpload>
  <Part>
    <PartNumber>1</PartNumber>
    <ETag>testetag1</ETag>
  </Part>
</CompleteMultipartUpload>`
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "testhost", req.Header.Get("Host"))
	assert.Equal(t, hashSHA256(expectXMLString), req.Header.Get("x-amz-content-sha256"))
	assert.Equal(t, "testAuthorization", req.Header.Get("Authorization"))
	assert.Equal(t, strconv.Itoa(len(expectXMLString)), req.Header.Get("Content-Length"))
}

func TestGenerateXML(t *testing.T) {
	etagMap := make(map[int]string)
	etagMap[1] = "testetag1"
	etagMap[2] = "testetag2"
	etagMap[3] = "testetag3"
	expectString := `<CompleteMultipartUpload>
  <Part>
    <PartNumber>1</PartNumber>
    <ETag>testetag1</ETag>
  </Part>
  <Part>
    <PartNumber>2</PartNumber>
    <ETag>testetag2</ETag>
  </Part>
  <Part>
    <PartNumber>3</PartNumber>
    <ETag>testetag3</ETag>
  </Part>
</CompleteMultipartUpload>`
	upload := &S3Upload{etagMapper: etagMap}
	actualString, _ := upload.generateXML()
	assert.Equal(t, expectString, actualString)
}

func TestDevideFile(t *testing.T) {
	upload := &S3Upload{}
	file, _ := os.Open("testdata/earth.jpg")
	upload.file = file
	upload.devideFile()
	assert.Equal(t, 4, len(upload.fileSlice))
}

func TestMutexMapInsert(t *testing.T) {
	etagMap := make(map[int]string)
	mutex := new(sync.Mutex)
	uploader := &S3Upload{
		etagMapper: etagMap,
		mutex:      mutex,
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		uploader.mutexMapInsert(1, "hoge")
		wg.Done()
	}()
	go func() {
		uploader.mutexMapInsert(2, "fuga")
		wg.Done()
	}()
	wg.Wait()

	assert.Equal(t, "hoge", uploader.etagMapper[1])
	assert.Equal(t, "fuga", uploader.etagMapper[2])
}
