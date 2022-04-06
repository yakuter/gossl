package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdHelp     = "help"
	CmdVerify   = "verify"
	CmdGenerate = "generate"
)

const (
	flagCAFile = "CAFile"
)

func main() {
	app := &cli.App{
		Name:     "GoSSL",
		Usage:    "Don't be afraid of SSL anymore",
		Commands: Commands(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:            CmdHelp,
			HelpName:        CmdHelp,
			Action:          CmdAction,
			Usage:           `displays help messages.`,
			Description:     `Display help messages.`,
			SkipFlagParsing: true,
			HideHelp:        true,
			HideHelpCommand: true,
		},
		{
			Name:        CmdVerify,
			HelpName:    CmdVerify,
			Action:      CmdAction,
			ArgsUsage:   `[cert file path]`,
			Usage:       `verifies certificate file.`,
			Description: `Verifies certificate file with provided CA file.`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  flagCAFile,
					Usage: "CA file path",
				},
			},
		},
	}
}

func CmdAction(c *cli.Context) error {
	if c.Err() != nil {
		return c.Err()
	}

	var err error
	switch cmd := c.Command.Name; cmd {
	case CmdHelp:
		if c.NArg() == 0 {
			err = cli.ShowAppHelp(c)
		} else {
			err = cli.ShowCommandHelp(c, c.Args().First())
		}
		if err != nil {
			return err
		}
	case CmdVerify:
		if err = ActionVerify(c); err != nil {
			return err
		}
	case CmdGenerate:
		if err = ActionGenerate(c); err != nil {
			return err
		}
	case "":
		if c.Args().Len() == 0 {
			break
		}
		fallthrough
	default:
		args := c.Args().Slice()
		if cmd == "" && len(args) > 0 {
			cmd = args[0]
			args = args[1:]
		}
		err = fmt.Errorf("command not found: '%s' arguments: %q", cmd, args)
	}

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionVerify(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("cert file is not provided")
	}

	// Check CA file if exist
	caFilePath := c.String(flagCAFile)
	caFileInfo, err := os.Lstat(caFilePath)
	if err != nil {
		return err
	}

	// Check cert file if exist
	certFilePath := c.Args().First()
	certFileInfo, err := os.Lstat(certFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("CA File info %s\n", caFileInfo.Name())
	fmt.Printf("Cert File info %s\n", certFileInfo.Name())

	fmt.Printf("Verify command args: %q\n", c.Args().Slice())
	return nil
}

func ActionGenerate(c *cli.Context) error {
	args := c.Args().Slice()
	fmt.Printf("Generate command args: %q\n", args)
	return nil
}
