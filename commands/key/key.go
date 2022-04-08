package key

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdKey = "key"
)

const (
	flagOut = "out"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdKey,
		HelpName:    CmdKey,
		Action:      Action,
		ArgsUsage:   `[numbits]`,
		Usage:       `generates RSA private key.`,
		Description: `Generates RSA private key with provided number of bits.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagOut,
			Usage:    "Output file name (optional)",
			Required: false,
		},
	}
}

func Action(c *cli.Context) error {
	log.Printf("Key command args: %q\n", c.Args().Slice())

	// Set numbits as int
	numbitsArg := c.Args().First()
	numbits, err := strconv.Atoi(numbitsArg)
	if err != nil {
		return errors.Wrapf(err, "failed to convert numbits %q to int error", numbitsArg)
	}

	// Set output
	var output *os.File = os.Stdout

	// If output file is provided, then create it and set as output
	if c.IsSet(flagOut) {
		outputFilePath := c.String(flagOut)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to create output file %q error", outputFilePath)
		}

		defer func() {
			if err := outputFile.Close(); err != nil {
				log.Printf("failed to close output file %q error %v", outputFilePath, err)
			}
		}()

		output = outputFile
	}

	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, numbits)
	if err != nil {
		return errors.Wrapf(err, "failed to generate rsa key error")
	}

	// Encode private key as PEM format
	privKeyPEM := bytes.NewBuffer(nil)
	err = pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err != nil {
		return errors.Wrapf(err, "failed to encode private key to pem error")
	}

	// Write private key to output
	_, err = output.WriteString(privKeyPEM.String())
	if err != nil {
		return errors.Wrapf(err, "failed to write output error")
	}

	log.Printf("Private key generated")

	return nil
}
