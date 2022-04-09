package key

import (
	"log"
	"os"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdKey = "key"

	flagOut  = "out"
	flagBits = "bits"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdKey,
		HelpName:    CmdKey,
		Action:      Action,
		ArgsUsage:   ` `,
		Usage:       `generates RSA private key.`,
		Description: `Generates RSA private key with provided number of bits.`,
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
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	// Write private key to output
	if err = os.WriteFile(outputFilePath, privateKeyBytes, 0o600); err != nil {
		log.Printf("Failed to write Private Key to file %s error: %v", output.Name(), err)
		return err
	}

	log.Printf("Private key generated")
	return nil
}
