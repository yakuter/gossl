package key

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdKey = "key"

	flagOut        = "out"
	flagBits       = "bits"
	flagWithPublic = "withpub"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdKey,
		HelpName:    CmdKey,
		Action:      Action,
		ArgsUsage:   ` `,
		Usage:       `generates RSA private and public key.`,
		Description: `Generates RSA private and public key with provided number of bits.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagOut,
			Usage:    "Output file path",
			Required: false,
		},
		&cli.UintFlag{
			Name:     flagBits,
			Usage:    "Number of bits",
			Required: true,
		},
		&cli.BoolFlag{
			Name:        flagWithPublic,
			Usage:       "Export public key with private key",
			Required:    false,
			DefaultText: "false",
		},
	}
}

func Action(c *cli.Context) error {
	// Generate private key
	privateKey, err := utils.GeneratePrivateKey(int(c.Uint(flagBits)))
	if err != nil {
		log.Printf("Failed to generate RSA Private Key with bit size %d error: %v", c.Uint(flagBits), err)
		return err
	}

	// Encode Private Key from RSA to PEM format
	privateKeyBytes := utils.PrivateKeyToPEM(privateKey)

	// Set output
	var (
		output             *os.File = os.Stdout
		outputPrivFilePath string   = output.Name()
		outputPubFilePath  string   = output.Name()
	)

	if c.IsSet(flagOut) {
		outputPrivFilePath = c.String(flagOut)
		outputPubFilePath = strings.TrimSuffix(outputPrivFilePath, filepath.Ext(outputPrivFilePath)) + ".pub"
	}

	// Write private key to output
	if err = os.WriteFile(outputPrivFilePath, privateKeyBytes, 0o600); err != nil {
		log.Printf("Failed to write Private Key to file %s error: %v", output.Name(), err)
		return err
	}

	log.Printf("Private key generated")

	// Export public key if flag is set
	if c.Bool(flagWithPublic) {
		publicKeyBytes := utils.PublicKeyToPEM(&privateKey.PublicKey)

		// Write public key to output
		if err = os.WriteFile(outputPubFilePath, publicKeyBytes, 0o600); err != nil {
			log.Printf("Failed to write Public Key to file %s error: %v", output.Name(), err)
			return err
		}
		log.Printf("Public key generated")
	}

	return nil
}
