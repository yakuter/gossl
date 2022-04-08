package verify

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"

	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdVerify = "verify"
)

const (
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
	log.Printf("Verify command args: %q\n", c.Args().Slice())

	// Read CA file
	caFilePath := c.String(flagCAFile)
	caFileBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read CA file %q content error", caFilePath)
	}

	// Read cert file
	certFilePath := c.String(flagCertFile)
	certFileBytes, err := os.ReadFile(certFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to read cert file %q content error", certFilePath)
	}

	// Generate an empty cert pool
	roots := x509.NewCertPool()

	// Append CA to cert pool
	if ok := roots.AppendCertsFromPEM(caFileBytes); !ok {
		err := fmt.Errorf("failed to append CA file %q to cert pool", caFilePath)
		return errors.Wrap(err, "error")
	}

	// Decode pem encoded cert file
	block, _ := pem.Decode(certFileBytes)
	if block == nil {
		err := fmt.Errorf("failed to decode PEM encoded cert file %q", certFilePath)
		return errors.Wrap(err, "error")
	}

	// Parse x509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.Wrapf(err, "failed to parse x509 certificate from cert file '%s' error", certFilePath)
	}

	// Set verification options
	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: x509.NewCertPool(),
	}

	// Check and add dns flag as DNSName
	if c.IsSet(flagDNS) {
		opts.DNSName = c.String(flagDNS)
	}

	// Verify certificate with verification options
	if _, err := cert.Verify(opts); err != nil {
		return errors.Wrapf(err, "failed to verify certificate CA file '%s', cert file '%s' error", caFilePath, certFilePath)
	}

	log.Printf("Certificate verification succeeds")

	return nil
}
