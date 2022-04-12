package ssh_copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/pkg/sftp"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

const (
	testUser = "testUser"
	testPass = "testPass"
	port     = "2022"
)

func TestSSHCopy(t *testing.T) {
	// Create an SSH server
	listener, err := net.Listen("tcp", ":"+port)
	require.NoError(t, err)

	defer listener.Close()

	// Use password authentication in SSH server
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == testUser && string(pass) == testPass {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	// A host key must be added to the server even password auth is used
	privateKey, err := utils.GeneratePrivateKey(1024)
	require.NoError(t, err)

	privateKeyBytes := utils.PrivateKeyToPEM(privateKey)

	private, err := ssh.ParsePrivateKey(privateKeyBytes)
	require.NoError(t, err)

	config.AddHostKey(private)

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleClient(nConn, config)
		}
	}()

	pubFile, err := os.CreateTemp(t.TempDir(), "pubkey")
	require.NoError(t, err)

	defer pubFile.Close()

	execName, err := os.Executable()
	require.NoError(t, err)

	testCases := []struct {
		name      string
		user      string
		pass      string
		shouldErr bool
	}{
		{
			name:      "valid sftp connection and write",
			user:      testUser,
			shouldErr: false,
		},
		{
			name:      "credentials error",
			user:      "wrongUser",
			shouldErr: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			testArgs := []string{
				execName, CmdSSHCopy,
				"--pubkey", pubFile.Name(),
				"--port", port,
				fmt.Sprintf("%s@localhost", tC.user),
			}

			app := &cli.App{
				Commands: []*cli.Command{
					Command(stubPasswordReader{Password: testPass}),
				},
			}

			if tC.shouldErr {
				require.Error(t, app.Run(testArgs))
			} else {
				require.NoError(t, app.Run(testArgs))

				// Make sure to delete the created local folder ".ssh"
				// in the current directory where this test is running
				// after a successful SSH copying
				currentDir, err := os.Getwd()
				require.NoError(t, err)
				sshFolderPath := filepath.Join(currentDir, ".ssh")

				defer os.RemoveAll(sshFolderPath)

				// Compare created public key with written authorized_keys
				// file after a successful SSH copying
				authorizedKeys := filepath.Join(sshFolderPath, "authorized_keys")
				authorizedKeysBytes, err := ioutil.ReadFile(authorizedKeys)
				require.NoError(t, err)

				tempPublicKey, err := ioutil.ReadAll(pubFile)
				require.NoError(t, err)

				require.Equal(t, authorizedKeysBytes, tempPublicKey)
			}
		})
	}
}

// handleClient handles an SFTP connection
func handleClient(nConn net.Conn, config *ssh.ServerConfig) error {
	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		return fmt.Errorf("failed to handshake: %v", err)
	}

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	newChannel := <-chans
	// Channels have a type, depending on the application level
	// protocol intended. In the case of an SFTP session, this is "subsystem"
	// with a payload string of "<length=4>sftp"
	if newChannel.ChannelType() != "session" {
		return fmt.Errorf("unknown channel type: %v", newChannel.ChannelType())
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		return fmt.Errorf("could not accept channel: %v", err)
	}

	// Sessions have out-of-band requests such as "shell",
	// "pty-req" and "env".  Here we handle only the
	// "subsystem" request.
	go func(in <-chan *ssh.Request) {
		for req := range in {
			ok := false
			switch req.Type {
			case "subsystem":
				if string(req.Payload[4:]) == "sftp" {
					ok = true
				}
			}
			req.Reply(ok, nil)
		}
	}(requests)

	server, err := sftp.NewServer(channel)
	if err != nil {
		return err
	}

	defer server.Close()

	err = server.Serve()
	if err != nil && err != io.EOF {
		return fmt.Errorf("sftp server returned with error: %v", err)
	}

	return nil
}

// stubPasswordReader is used to mock the password input given
type stubPasswordReader struct {
	Password string
}

func (pr stubPasswordReader) ReadPassword() (string, error) {
	return pr.Password, nil
}
