package uploader

import "testing"

func TestInitialMultipartUpload(t *testing.T) {
	upload, _ := New("s3go-cli-test", "test")
	upload.InitialMultipartUpload()
}
