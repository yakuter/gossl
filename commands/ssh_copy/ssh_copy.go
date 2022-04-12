package ssh_copy

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/sftp"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

const (
	CmdSSHCopy = "ssh-copy"

	flagPubkey = "pubkey"
	flagPort   = "port"
)

func Command(reader passwordReader) *cli.Command {
	return &cli.Command{
		Name:        CmdSSHCopy,
		HelpName:    CmdSSHCopy,
		Action:      Action(reader),
		ArgsUsage:   `[remote-user@remote-ip]`,
		Usage:       `copy SSH public key to remote server.`,
		Description: `Copy SSH public key to authorized_keys in remote SSH server.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        flagPubkey,
			Usage:       "Output file path",
			Required:    false,
			DefaultText: "eg, /home/user/.ssh/id_rsa.pub",
		},
		&cli.UintFlag{
			Name:        flagPort,
			Usage:       "SSH server connection port",
			Required:    false,
			DefaultText: "eg, 22",
			Value:       22,
		},
	}
}

func Action(reader passwordReader) func(*cli.Context) error {
	return func(c *cli.Context) error {
		// Read public key from file
		pubKey, err := os.ReadFile(c.String(flagPubkey))
		if err != nil {
			log.Printf("Failed to read public key error: %v", err)
			return err
		}

		// Parse remote username and hostname/ip
		remote := c.Args().First()
		user, host, found := strings.Cut(remote, "@")
		if !found {
			err := fmt.Errorf("failed to parse remote user and hostname from %s", remote)
			log.Printf("%v", err)
			return err
		}

		// Get password from user
		fmt.Printf("Password: ")
		pwd, err := reader.ReadPassword()
		if err != nil {
			log.Printf("failed to read inputs %v", err)
			return err
		}

		// Connect to remote SSH server with SFTP
		client, err := connectSFTP(host, user, pwd, int(c.Uint(flagPort)))
		if err != nil {
			log.Printf("Failed to connect SSH server error: %v", err)
			return err
		}
		defer func() {
			if err = client.Close(); err != nil {
				log.Printf("Failed to close SSH connection error: %v", err)
			}
		}()

		// We expect working directory is SSH user's home directory
		workdir, err := client.Getwd()
		if err != nil {
			log.Printf("Failed to get working directory error: %v", err)
			return err
		}

		// Create .ssh and its parent folders if not exist
		remoteDir := filepath.Join(workdir, ".ssh")
		if err = client.MkdirAll(remoteDir); err != nil {
			log.Printf("Failed to create remote dir %s error: %v", remoteDir, err)
			return err
		}

		// Open or create authorized_keys to append public key
		remotePath := filepath.Join(remoteDir, "authorized_keys")
		file, err := client.OpenFile(remotePath, os.O_RDWR|os.O_APPEND|os.O_CREATE)
		if err != nil {
			log.Printf("Failed to open file authorized_keys error: %v", err)
			return err
		}

		defer func() {
			if err = file.Close(); err != nil {
				log.Printf("Failed to close remote authorized_keys file error: %v", err)
			}
		}()

		if _, err = file.Write(pubKey); err != nil {
			log.Printf("Failed to write public key to remote authorized_keys error: %v", err)
			return err
		}

		log.Printf("SSH Public Key added to remote server")
		return nil
	}
}

// connectSFTP creates ssh config, tries to connect (dial) SSH server and
// creates new client with connection
func connectSFTP(host, username, password string, port int) (*sftp.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("Failed to dial error: %v ", err)
		return nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Printf("Failed to create new SFTP client Error: %v", err)
		return nil, err
	}

	return client, nil
}

type passwordReader interface {
	ReadPassword() (string, error)
}

type StdinPasswordReader struct{}

func (pr StdinPasswordReader) ReadPassword() (string, error) {
	pwd, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(pwd), nil
}
