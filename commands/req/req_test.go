package req_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/yakuter/gossl/commands/key"
	"github.com/yakuter/gossl/commands/req"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestReq(t *testing.T) {
	execName, err := os.Executable()
	require.NoError(t, err)

	tempDir := t.TempDir()
	outFile := filepath.Join(tempDir, "test.cert")
	testKey := filepath.Join(tempDir, "test.key")

	keyApp := &cli.App{Commands: []*cli.Command{key.Command()}}
	err = keyApp.Run([]string{execName, key.CmdKey, "-out", testKey, "-bits", "2048"})
	require.NoError(t, err)

	testCases := []struct {
		name      string
		email     string
		fqdn      string
		key       string
		out       string
		days      int
		serial    int
		isCA      bool
		shouldErr bool
	}{
		{
			name:      "valid CA",
			fqdn:      "localhost",
			key:       testKey,
			out:       outFile,
			days:      365,
			serial:    123456,
			isCA:      true,
			shouldErr: false,
		},
		{
			name:      "valid CSR",
			fqdn:      "localhost",
			email:     "john@doe.com",
			key:       testKey,
			out:       outFile,
			days:      365,
			serial:    123456,
			isCA:      false,
			shouldErr: false,
		},
		{
			name:      "empty email CSR error",
			fqdn:      "localhost",
			email:     "",
			key:       testKey,
			out:       outFile,
			days:      365,
			serial:    123456,
			isCA:      false,
			shouldErr: true,
		},
		{
			name:      "empty FQDN error",
			fqdn:      "",
			key:       testKey,
			out:       outFile,
			days:      365,
			serial:    123456,
			isCA:      true,
			shouldErr: true,
		},
		{
			name:      "key error",
			fqdn:      "localhost",
			key:       "",
			out:       outFile,
			days:      365,
			serial:    123456,
			isCA:      true,
			shouldErr: true,
		},
		{
			name:      "output file error",
			fqdn:      "localhost",
			key:       testKey,
			out:       "",
			days:      365,
			serial:    123456,
			isCA:      true,
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, req.CmdCert,
				"--key", tC.key,
				"--out", tC.out,
				"--days", strconv.Itoa(tC.days),
				"--serial", strconv.Itoa(tC.serial),
			}
			if tC.isCA {
				testArgs = append(testArgs, "--isCA")
			}

			var stdin bytes.Buffer
			stdin.Write([]byte(tC.fqdn + "\n" + tC.email + "\na\na\na\na\na\na\na"))

			app := &cli.App{
				Commands: []*cli.Command{
					req.Command(&stdin),
				},
			}

			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}
