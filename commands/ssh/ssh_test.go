package ssh_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yakuter/gossl/commands/ssh"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestKey(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			ssh.Command(),
		},
	}

	execName, err := os.Executable()
	require.NoError(t, err)

	tempDir := t.TempDir()
	outFilePath := filepath.Join(tempDir, "id_rsa")

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
			testArgs := []string{execName, ssh.CmdSSH, "-out", tC.out, "-bits", tC.numbits}
			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
				require.FileExists(t, outFilePath+".pub")
			}
		})
	}
}
