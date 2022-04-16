package main

import (
	"io"
	"log"
	"os"

	"github.com/yakuter/gossl/commands/cert"
	"github.com/yakuter/gossl/commands/help"
	"github.com/yakuter/gossl/commands/info"
	"github.com/yakuter/gossl/commands/key"
	"github.com/yakuter/gossl/commands/ssh"
	"github.com/yakuter/gossl/commands/ssh_copy"
	"github.com/yakuter/gossl/commands/verify"

	"github.com/urfave/cli/v2"
)

var Version = "v0.1.6"

func main() {
	app := &cli.App{
		Name:     "GoSSL",
		Usage:    "Don't be afraid of SSL anymore",
		Commands: Commands(os.Stdin),
		Version:  Version,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Commands(reader io.Reader) []*cli.Command {
	return []*cli.Command{
		help.Command(),
		key.Command(),
		cert.Command(reader),
		info.Command(),
		verify.Command(),
		ssh.Command(),
		ssh_copy.Command(ssh_copy.StdinPasswordReader{}),
	}
}
