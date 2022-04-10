package ssh_upload

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"github.com/urfave/cli/v2"
	"github.com/yakuter/gossl/pkg/utils"
	"golang.org/x/crypto/ssh"
)

const (
	CmdSSHUpload = "ssh-upload"

	flagPubkey = "pubkey"
	flagPort   = "port"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdSSHUpload,
		HelpName:    CmdSSHUpload,
		Action:      Action,
		ArgsUsage:   ` `,
		Usage:       `generates RSA SSH key pair.`,
		Description: `Generates RSA SSH key pair private and public key with provided number of bits.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        flagPubkey,
			Usage:       "Output file path",
			Required:    false,
			DefaultText: "eg, ./id_rsa.pub",
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

func Action(c *cli.Context) error {
	// Read public key to send remote SSH server
	pubKey, err := os.ReadFile(c.String(flagPubkey))
	if err != nil {
		log.Printf("Failed to read public key error: %v", err)
		return err
	}

	// Parser remote username and hostname/ip
	remote := c.Args().First()
	user, host, found := strings.Cut(remote, "@")
	if !found {
		err := fmt.Errorf("failed to parse remote user and hostname from %s", remote)
		log.Printf("%v", err)
		return err
	}

	// Get Password from user
	answers, err := utils.ReadInputs([]string{"Password"})
	if err != nil {
		log.Printf("failed to read inputs %v", err)
		return err
	}

	// Connect to remote SSH server
	client, err := connectSFTP(host, user, answers[0], int(c.Uint(flagPort)))
	if err != nil {
		log.Printf("Failed to connect SSH server error: %v", err)
		return err
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Printf("Failed to close SSH connection error: %v", err)
		}
	}()

	workdir, err := client.Getwd()
	if err != nil {
		log.Printf("Failed to get working directory error: %v", err)
		return err
	}

	remoteDir := filepath.Join(workdir, ".ssh")

	if err = client.MkdirAll(remoteDir); err != nil {
		log.Printf("Failed to create remote dir %s error: %v", remoteDir, err)
		return err
	}

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
