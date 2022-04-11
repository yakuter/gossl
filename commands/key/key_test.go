package key_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yakuter/gossl/commands/key"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestKey(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			key.Command(),
		},
	}

	execName, err := os.Executable()
	require.NoError(t, err)

	tempDir := t.TempDir()
	outFilePath := filepath.Join(tempDir, "private.key")

	testCases := []struct {
		name      string
		out       string
		filePath  string
		numbits   string
		shouldErr bool
	}{
		{
			name:      "valid private key",
			out:       outFilePath,
			numbits:   "1024",
			shouldErr: false,
		},
		{
			name:      "out file not found error",
			out:       "",
			numbits:   "1024",
			shouldErr: true,
		},
		{
			name:      "numbits not integer",
			out:       outFilePath,
			numbits:   "not-number",
			shouldErr: true,
		},
		{
			name:      "numbits not integer",
			out:       outFilePath,
			numbits:   "",
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, key.CmdKey, "-out", tC.out, "-bits", tC.numbits, "-withpub"}
			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}
