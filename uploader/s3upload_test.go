package uploader

import (
	"testing"

	"github.com/hikaru7719/s3go/signature"
)

func TestInitialMultipartUpload(t *testing.T) {
	sig := signature.New()
	upload, _ := New("s3go-cli-test", "test", sig)
	upload.InitialMultipartUpload()
}
