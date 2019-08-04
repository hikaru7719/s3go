package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		fmt.Println("cli execute")
		return nil
	}
}
