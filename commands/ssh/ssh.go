package ssh

import (
	"log"
	"os"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdSSH = "ssh"

	flagOut  = "out"
	flagBits = "bits"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdSSH,
		HelpName:    CmdSSH,
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
			Name:        flagOut,
			Usage:       "Output file path",
			Required:    false,
			Value:       "./id_rsa",
			DefaultText: "./id_rsa",
		},
		&cli.UintFlag{
			Name:     flagBits,
			Usage:    "Number of bits",
			Required: true,
		},
	}
}

func Action(c *cli.Context) error {
	// Generate Private Key
	privateKey, err := utils.GeneratePrivateKey(int(c.Uint(flagBits)))
	if err != nil {
		log.Printf("Failed to generate RSA Private Key with bit size %d error: %v", c.Uint(flagBits), err)
		return err
	}

	// Generate Public Key from Private Key
	publicKeyBytes, err := utils.GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Printf("Failed to generate Public Key error: %v", err)
		return err
	}

	// Encode Private Key from RSA to PEM format
	privateKeyBytes := utils.PrivateKeyToPEM(privateKey)

	// Write private key to file
	privateKeyFilePath := c.String(flagOut)
	if err = os.WriteFile(privateKeyFilePath, privateKeyBytes, 0o600); err != nil {
		log.Printf("Failed to write Private Key to file %s error: %v", privateKeyFilePath, err)
		return err
	}

	// Write public key to file
	publicKeyFilePath := privateKeyFilePath + ".pub"
	if err = os.WriteFile(publicKeyFilePath, publicKeyBytes, 0o600); err != nil {
		log.Printf("Failed to write Public Key to file %s error: %v", publicKeyFilePath, err)
		return err
	}

	log.Printf("SSH key pair generated")
	return nil
}
