package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hikaru7719/s3go/signature"
	"github.com/hikaru7719/s3go/uploader"
	"github.com/urfave/cli"
)

func main() {
	app := New()
	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// New create cli App
func New() *cli.App {
	app := cli.NewApp()
	app.Name = "s3go"
	app.Usage = "Upload some file to AWS S3"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "`File` to upload to S3",
		},
		cli.StringFlag{
			Name:  "bucket, b",
			Value: "",
			Usage: "`S3 bucket Name` to upload files",
		},
	}

	app.Action = func(c *cli.Context) error {
		var file, bucket string
		if f := c.String("file"); f != "" {
			file = f
		}

		if b := c.String("bucket"); b != "" {
			bucket = b
		}

		fmt.Println(file, bucket)
		sign := signature.New()
		uploader, err := uploader.New(bucket, file, sign)
		if err != nil {
			return nil
		}

		return uploader.Run()
	}
	return app
}
