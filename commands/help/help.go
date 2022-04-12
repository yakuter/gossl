package help

import (
	"github.com/urfave/cli/v2"
)

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

func Action(c *cli.Context) error {
	var err error
	if c.NArg() == 0 {
		err = cli.ShowAppHelp(c)
	} else {
		err = cli.ShowCommandHelp(c, c.Args().First())
	}
	return err
}
