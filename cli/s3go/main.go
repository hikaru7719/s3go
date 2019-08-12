package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := New()
	app.Run(os.Args)
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
			Usage: "file to upload to S3",
		},
		cli.StringFlag{
			Name:  "bucket, b",
			Value: "",
			Usage: "S3 bucket Name to be upload files",
		},
		cli.StringFlag{
			Name:  "key, k",
			Value: "",
			Usage: "AWS access key",
		},
		cli.StringFlag{
			Name:  "secret, s",
			Value: "",
			Usage: "AWS secret",
		},
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "AWS region",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("cli execute")
		return nil
	}
	return app
}
