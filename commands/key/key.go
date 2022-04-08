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

	"github.com/urfave/cli/v2"
)

const (
	CmdKey = "key"

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
	numbitsArg := c.Args().First()

	log.Printf("Generating private key with number of bits %q", numbitsArg)

	// Set numbits as int
	numbits, err := strconv.Atoi(numbitsArg)
	if err != nil {
		log.Printf("Failed to convert numbits %q to int error: %v", numbitsArg, err)
		return err
	}

	// Set output
	var output *os.File = os.Stdout

	// If output file is provided, then create it and set as output
	if c.IsSet(flagOut) {
		outputFilePath := c.String(flagOut)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Printf("Failed to create output file %q error: %v", outputFilePath, err)
			return err
		}

		defer func() {
			if err := outputFile.Close(); err != nil {
				log.Printf("Failed to close output file %q error: %v", outputFilePath, err)
			}
		}()

		output = outputFile
	}

	// Generate private key
	privKey, err := rsa.GenerateKey(rand.Reader, numbits)
	if err != nil {
		log.Printf("Failed to generate rsa key error: %v", err)
		return err
	}

	// Encode private key as PEM format
	privKeyPEM := bytes.NewBuffer(nil)
	err = pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err != nil {
		log.Printf("Failed to encode private key as PEM error: %v", err)
		return err
	}

	// Write private key to output
	_, err = output.WriteString(privKeyPEM.String())
	if err != nil {
		log.Printf("Failed to write private key to output error: %v", err)
		return err
	}

	log.Printf("Private key generated")
	return nil
}
