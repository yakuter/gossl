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
		argFile               = "../../testdata/server-cert.pem"
		flagURL               = "google.com"
		flagURLwithScheme     = "https://google.com"
		flagURLwithPort       = "google.com:443"
		flagURLwithSchemePort = "https://google.com:443"
	)

	testCases := []struct {
		name      string
		arg       string
		url       string
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
			url:       flagURL,
			shouldErr: false,
		},
		{
			name:      "wrong output with url",
			url:       flagURL,
			out:       "/wrong-out",
			shouldErr: true,
		},
		{
			name:      "valid URL with output",
			url:       flagURL,
			out:       outFilePath,
			shouldErr: false,
		},
		{
			name:      "valid URL with scheme",
			url:       flagURLwithScheme,
			shouldErr: false,
		},
		{
			name:      "valid URL with port",
			url:       flagURLwithPort,
			shouldErr: false,
		},
		{
			name:      "valid URL with scheme and port",
			url:       flagURLwithSchemePort,
			shouldErr: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, info.CmdInfo}
			if tC.out != "" {
				testArgs = append(testArgs, "--out", tC.out)
			}
			if tC.url != "" {
				testArgs = append(testArgs, "--url", tC.url)
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
