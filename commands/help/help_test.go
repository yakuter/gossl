package help_test

import (
	"os"
	"testing"

	"github.com/yakuter/gossl/commands/help"
	"github.com/yakuter/gossl/commands/key"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestHelp(t *testing.T) {
	execName, err := os.Executable()
	require.NoError(t, err)

	app := &cli.App{
		Commands: []*cli.Command{
			help.Command(),
			key.Command(),
		},
	}

	testArgs := []string{execName, "help"}
	require.NoError(t, app.Run(testArgs))

	testArgs = []string{execName, "help", "key"}
	require.NoError(t, app.Run(testArgs))
}
