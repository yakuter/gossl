package help

import (
	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdHelp = "help"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            CmdHelp,
		HelpName:        CmdHelp,
		Action:          Action,
		Usage:           `displays help messages.`,
		Description:     `Display help messages.`,
		SkipFlagParsing: true,
		HideHelp:        true,
		HideHelpCommand: true,
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		// ...
	}
}

func Action(c *cli.Context) error {
	var err error
	if c.NArg() == 0 {
		err = cli.ShowAppHelp(c)
	} else {
		err = cli.ShowCommandHelp(c, c.Args().First())
	}
	if err != nil {
		return err
	}

	return nil
}
