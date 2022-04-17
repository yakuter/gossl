package info_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yakuter/gossl/commands/info"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestInfo(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			info.Command(),
		},
	}

	execName, err := os.Executable()
	require.NoError(t, err)

	tempDir := t.TempDir()
	outFilePath := filepath.Join(tempDir, "test-file-*")

	arg := "../../testdata/server-cert.pem"

	testCases := []struct {
		name      string
		arg       string
		out       string
		shouldErr bool
	}{
		{
			name:      "valid cert",
			arg:       arg,
			shouldErr: false,
		},
		{
			name:      "no ",
			shouldErr: true,
		},
		{
			name:      "wrong cert",
			arg:       "wrong-arg",
			shouldErr: true,
		},
		{
			name:      "valid cert",
			arg:       arg,
			out:       outFilePath,
			shouldErr: false,
		},
		{
			name:      "wrong output",
			arg:       arg,
			out:       "/wrong-out",
			shouldErr: true,
		},
		{
			name:      "valid URL",
			arg:       "https://www.google.com",
			shouldErr: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, info.CmdInfo}
			if tC.out != "" {
				testArgs = append(testArgs, "--out", tC.out)
			}
			if tC.arg != "" {
				testArgs = append(testArgs, tC.arg)
			}

			err = app.Run(testArgs)
			if tC.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
