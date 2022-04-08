package verify

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdVerify = "verify"

	flagCAFile   = "cafile"
	flagCertFile = "certfile"
	flagDNS      = "dns"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdVerify,
		HelpName:    CmdVerify,
		Action:      Action,
		ArgsUsage:   `[cert file path]`,
		Usage:       `verifies certificate file.`,
		Description: `Verifies certificate file with provided CA file.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagCAFile,
			Usage:    "CA file path (required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagCertFile,
			Usage:    "Cert file path (required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagDNS,
			Usage:    "DNS name or IP (optional)",
			Required: false,
		},
	}
}

func Action(c *cli.Context) error {
	var (
		caFilePath   = c.String(flagCAFile)
		certFilePath = c.String(flagCertFile)
	)

	log.Printf("Verification started with CA file %q and cert file %q\n", caFilePath, certFilePath)

	// Read CA file
	caFileBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Printf("Failed to read CA file %q error: %v", caFilePath, err)
		return err
	}

	// Read cert file
	certFileBytes, err := os.ReadFile(certFilePath)
	if err != nil {
		log.Printf("Failed to read cert file %q error: %v", certFilePath, err)
		return err
	}

	// Generate an empty cert pool
	roots := x509.NewCertPool()

	// Append CA to cert pool
	if ok := roots.AppendCertsFromPEM(caFileBytes); !ok {
		err = fmt.Errorf("Failed to append CA from %q to cert pool", caFilePath)
		log.Printf("%v", err)
		return err
	}

	// Decode PEM encoded cert file
	block, _ := pem.Decode(certFileBytes)
	if block == nil {
		err = errors.New("block is nil")
		log.Printf("Failed to decode PEM encoded cert file %q error: %v", certFilePath, err)
		return err
	}

	// Parse x509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Printf("Failed to parse x509 certificate from cert file %q error: %v", certFilePath, err)
		return err
	}

	if err = verify(c, cert, roots); err != nil {
		log.Printf("Failed to verify CA and cert error: %v", err)
		return err
	}

	log.Printf("Certificate verification succeeds")
	return nil
}

func verify(c *cli.Context, cert *x509.Certificate, roots *x509.CertPool) error {
	// Set verification options
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
	}

	// Add dns flag as DNSName if set
	if c.IsSet(flagDNS) {
		opts.DNSName = c.String(flagDNS)
	}

	// Verify certificate with verification options
	if _, err := cert.Verify(opts); err != nil {
		log.Printf("Failed to verify certificate error: %v", err)
		return err
	}

	return nil
}
