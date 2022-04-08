package verify_test

import (
	"os"
	"testing"

	"github.com/yakuter/gossl/commands/verify"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
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
		dns          = "127.0.0.1"
		caFilePath   = "../../testdata/ca-cert.pem"
		certFilePath = "../../testdata/server-cert.pem"
	)

	testCases := []struct {
		name      string
		cafile    string
		certfile  string
		dns       string
		shouldErr bool
	}{
		{
			name:      "valid cert and ca",
			cafile:    "../../testdata/ca-cert.pem",
			certfile:  "../../testdata/server-cert.pem",
			dns:       dns,
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
			dns:       dns,
			shouldErr: true,
		},
		{
			name:      "cert file error",
			cafile:    caFilePath,
			certfile:  "wrong-file",
			dns:       dns,
			shouldErr: true,
		},
		{
			name:      "dns error",
			cafile:    caFilePath,
			certfile:  certFilePath,
			dns:       "wrong.dns.com",
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{execName, verify.CmdVerify, "-dns", tC.dns, "-cafile", tC.cafile, "-certfile", tC.certfile}
			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}
