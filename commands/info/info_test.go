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

	const (
		argFile = "../../testdata/server-cert.pem"
		argURL  = "https://www.google.com"
	)

	testCases := []struct {
		name      string
		arg       string
		out       string
		shouldErr bool
	}{
		{
			name:      "valid cert file",
			arg:       argFile,
			shouldErr: false,
		},
		{
			name:      "no ",
			shouldErr: true,
		},
		{
			name:      "wrong cert file or invalid URL",
			arg:       "wrong-arg",
			shouldErr: true,
		},
		{
			name:      "valid cert file with output",
			arg:       argFile,
			out:       outFilePath,
			shouldErr: false,
		},
		{
			name:      "wrong output",
			arg:       argFile,
			out:       "/wrong-out",
			shouldErr: true,
		},
		{
			name:      "valid URL",
			arg:       argURL,
			shouldErr: false,
		},
		{
			name:      "wrong output with url",
			arg:       argURL,
			out:       "/wrong-out",
			shouldErr: true,
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
