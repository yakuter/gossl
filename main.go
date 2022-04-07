package main

import (
	"fmt"
	"log"
	"os"

	"github.com/yakuter/gossl/commands/help"
	"github.com/yakuter/gossl/commands/key"
	"github.com/yakuter/gossl/commands/verify"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "GoSSL",
		Usage:    "Don't be afraid of SSL anymore",
		Commands: Commands(),
	}

	fmt.Println(os.Args)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Commands() []*cli.Command {
	return []*cli.Command{
		help.Command(),
		verify.Command(),
		key.Command(),
	}
}
