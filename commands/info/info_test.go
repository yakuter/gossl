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
		certFile              = "../../testdata/server-cert.pem"
		csrFile               = "../../testdata/server-req.pem"
		flagURL               = "google.com"
		flagURLwithScheme     = "https://google.com"
		flagURLwithPort       = "google.com:443"
		flagURLwithSchemePort = "https://google.com:443"
	)

	testCases := []struct {
		name      string
		cert      string
		csr       string
		url       string
		out       string
		shouldErr bool
	}{
		{
			name:      "valid cert file",
			cert:      certFile,
			shouldErr: false,
		},
		{
			name:      "valid csr file",
			csr:       csrFile,
			shouldErr: false,
		},
		{
			name:      "no ",
			shouldErr: true,
		},
		{
			name:      "wrong cert file",
			cert:      "wrong-arg",
			shouldErr: true,
		},
		{
			name:      "wrong cert file",
			cert:      csrFile,
			shouldErr: true,
		},
		{
			name:      "wrong csr file",
			csr:       "wrong-arg",
			shouldErr: true,
		},
		{
			name:      "wrong csr file",
			csr:       certFile,
			shouldErr: true,
		},
		{
			name:      "valid cert file with output",
			cert:      certFile,
			out:       outFilePath,
			shouldErr: false,
		},
		{
			name:      "valid csr file with output",
			csr:       csrFile,
			out:       outFilePath,
			shouldErr: false,
		},
		{
			name:      "wrong output with cert",
			cert:      certFile,
			out:       "/wrong-out",
			shouldErr: true,
		},
		{
			name:      "wrong output with csr",
			csr:       csrFile,
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
			if tC.cert != "" {
				testArgs = append(testArgs, "--cert", tC.cert)
			}
			if tC.csr != "" {
				testArgs = append(testArgs, "--csr", tC.csr)
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
