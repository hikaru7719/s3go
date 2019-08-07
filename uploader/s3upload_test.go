package uploader

import (
	"os"
	"testing"

	"github.com/hikaru7719/s3go/signature"
	"github.com/stretchr/testify/assert"
)

func TestInitialMultipartUpload(t *testing.T) {
	sig := signature.New()
	upload, _ := New("s3go-cli-test", "test", sig)
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

func TestPutMultipartObject(t *testing.T) {
	sig := signature.New()
	upload, _ := New("s3go-cli-test", "test", sig)
	file, _ := os.Open("earth.jpg")
	upload.file = file
	upload.InitialMultipartUpload()
	upload.PutMaltiPartObject()
}

func TestNewUploadRequest(t *testing.T) {
	upload := &S3Upload{
		host:       "testhost",
		bucketName: "testbucket",
		objectName: "testObject",
		signature:  &mockAuth{},
	}
	file, _ := os.Open("earth.jpg")
	upload.file = file
	req, _ := upload.newUploaderRequest()
	assert.Equal(t, "hoge", req.Header.Get("content-length"))
}
