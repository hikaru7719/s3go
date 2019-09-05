package benchmark

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func BenchmarkS3goCLI(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./s3go", "-b", os.Getenv("AWS_S3_BUCKET_NAME"), "-f", "../testdata/earth.jpg")
		err := cmd.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAWSCLI(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("aws", "s3", "cp", "../testdata/earth2.jpg", fmt.Sprintf("s3://%s", os.Getenv("AWS_S3_BUCKET_NAME")))
		err := cmd.Run()
		if err != nil {
			b.Fatal(err)
		}
	}
}
