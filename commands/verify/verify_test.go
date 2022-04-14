package verify_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	emptyFile, err := os.CreateTemp(t.TempDir(), "empty-file-*")
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
			name:      "empty cert file error",
			cafile:    caFilePath,
			certfile:  emptyFile.Name(),
			dns:       dns,
			shouldErr: true,
		},
		{
			name:      "empty ca file error",
			cafile:    emptyFile.Name(),
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
			testArgs := []string{
				execName, verify.CmdVerify,
				"-dns", tC.dns,
				"-cafile", tC.cafile,
				"-certfile", tC.certfile,
			}
			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}

func TestVerifyURL(t *testing.T) {
	const (
		// Certificate paths
		serverCert = "../../testdata/server-cert.pem"
		serverKey  = "../../testdata/server-key.pem"
		caCert     = "../../testdata/ca-cert.pem"
		caCert2    = "../../testdata/ca-cert-2.pem"
	)

	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello!")
		}),
	)

	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	require.NoError(t, err)

	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	ts.StartTLS()
	defer ts.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			verify.Command(),
		},
	}

	execName, err := os.Executable()
	require.NoError(t, err)

	testCases := []struct {
		name      string
		cafile    string
		shouldErr bool
	}{
		{
			name:      "valid cert and ca",
			cafile:    caCert,
			shouldErr: false,
		},
		{
			name:      "bad certificate",
			cafile:    caCert2,
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{
				execName, verify.CmdVerify,
				"--cafile", tC.cafile,
				"--url", ts.URL,
			}

			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))
			}
		})
	}
}
