package verify_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"github.com/yakuter/gossl/commands/verify"
)

func TestVerify(t *testing.T) {
	app := &cli.App{
		Commands: []*cli.Command{
			verify.Command(),
		},
	}

	execName, err := os.Executable()
	require.NoError(t, err)

	var (
		hostname     = "127.0.0.1"
		caFilePath   = "../../testdata/ca-cert.pem"
		certFilePath = "../../testdata/server-cert.pem"
	)

	testCases := []struct {
		name      string
		cafile    string
		certfile  string
		hostname  string
		shouldErr bool
	}{
		{
			name:      "valid cert and ca",
			cafile:    "../../testdata/ca-cert.pem",
			certfile:  "../../testdata/server-cert.pem",
			hostname:  hostname,
			shouldErr: false,
		},
		{
			name:      "argument error",
			certfile:  "",
			shouldErr: true,
		},
		{
			name:      "ca file error",
			cafile:    "wrong-file",
			certfile:  certFilePath,
			hostname:  hostname,
			shouldErr: true,
		},
		{
			name:      "cert file error",
			cafile:    caFilePath,
			certfile:  "wrong-file",
			hostname:  hostname,
			shouldErr: true,
		},
		{
			name:      "hostname error",
			cafile:    caFilePath,
			certfile:  certFilePath,
			hostname:  "wrong.hostname.com",
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, verify.CmdVerify, "-hostname", tC.hostname, "-cafile", tC.cafile, tC.certfile}
			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}

func testFileWithContent(t *testing.T, tempdir, content string) string {
	file, err := os.CreateTemp(tempdir, "test-file-*")
	require.NoError(t, err)

	_, err = file.WriteString(content)
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	return file.Name()
}
