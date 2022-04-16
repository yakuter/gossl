package verify

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdVerify = "verify"

	flagCAFile   = "cafile"
	flagCertFile = "certfile"
	flagDNS      = "dns"
	flagURL      = "url"
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
			Usage:    "Cert file path to verify with CA (optional)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     flagDNS,
			Usage:    "DNS name or IP to verify with cert file and CA (optional)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     flagURL,
			Usage:    "URL to verify with CA (optional)",
			Required: false,
		},
	}
}

func Action(c *cli.Context) error {
	// At least one of the flag is required
	if !c.IsSet(flagURL) && !c.IsSet(flagCertFile) {
		return errors.New("Please provide url or certfile flag")
	}

	if c.IsSet(flagURL) && c.IsSet(flagDNS) {
		return errors.New("DNS flag is only allowed to be used with certfile flag")
	}

	// Generate new cert pool with CA file
	roots, err := rootCAs(c.String(flagCAFile))
	if err != nil {
		log.Printf("Failed to get root CAs error: %v", err)
		return err
	}

	// Verify cert file
	if c.IsSet(flagCertFile) {
		cert, err := utils.CertFromFile(c.String(flagCertFile))
		if err != nil {
			log.Printf("Failed to get cert from file %s CAs error: %v", c.String(flagCertFile), err)
			return err
		}

		if err = verifyCertWithCA(c, cert, roots); err != nil {
			log.Printf("Failed to verify CA and cert error: %v", err)
			return err
		}
	}

	// Verify URL
	if c.IsSet(flagURL) {
		if err = verifyURLWithCA(c, c.String(flagURL), roots); err != nil {
			log.Printf("Failed to verify CA and URL error: %v", err)
			return err
		}
	}

	log.Printf("Certificate verification succeeds")
	return nil
}

func rootCAs(caFilePath string) (*x509.CertPool, error) {
	// Read CA file
	caFileBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Printf("Failed to read CA file %q error: %v", caFilePath, err)
		return nil, err
	}

	// Generate an empty cert pool
	roots := x509.NewCertPool()

	// Append CA to cert pool
	if ok := roots.AppendCertsFromPEM(caFileBytes); !ok {
		err = fmt.Errorf("Failed to append CA from %q to cert pool", caFilePath)
		log.Printf("%v", err)
		return nil, err
	}

	return roots, nil
}

func verifyCertWithCA(c *cli.Context, cert *x509.Certificate, roots *x509.CertPool) error {
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

func verifyURLWithCA(c *cli.Context, url string, roots *x509.CertPool) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            roots,
			InsecureSkipVerify: false,
		},
	}

	client := &http.Client{Transport: tr}

	_, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to send Get request to URL %s error: %v", url, err)
		return err
	}

	return nil
}
