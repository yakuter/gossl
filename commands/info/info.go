package info

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/grantae/certinfo"
	"github.com/urfave/cli/v2"
)

const (
	CmdInfo = "info"

	flagOut  = "out"
	flagURL  = "url"
	flagCert = "cert"
	flagCSR  = "csr"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdInfo,
		HelpName:    CmdInfo,
		Action:      Action,
		ArgsUsage:   `[cert file path]`,
		Usage:       `displays information about certificate.`,
		Description: `Displays information about x509 certificate.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        flagOut,
			Usage:       "Output file name (optional)",
			DefaultText: "eg, ./cert.pem",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        flagURL,
			Usage:       "URL to get certificate details from (optional)",
			DefaultText: "eg, google.com",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        flagCert,
			Usage:       "x509 Certificate to get details from (optional)",
			DefaultText: "eg, server.crt",
			Required:    false,
		},
		&cli.StringFlag{
			Name:        flagCSR,
			Usage:       "x509 Certificate Request (CSR) to get details from (optional)",
			DefaultText: "eg, server.csr",
			Required:    false,
		},
	}
}

func Action(c *cli.Context) error {
	if !c.IsSet(flagURL) && !c.IsSet(flagCert) && !c.IsSet(flagCSR) {
		err := errors.New("no flag provided")
		log.Printf("Failed to get cert resource error: %v", err)
		return err
	}

	// Set output
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	var (
		err    error
		result string
		cert   *x509.Certificate
		csr    *x509.CertificateRequest
	)

	if c.IsSet(flagURL) {
		u := c.String(flagURL)
		cert, err = readX509FromDomain(u)
		if err != nil {
			log.Printf("failed to get cert details from URL %q error: %v", u, err)
			return err
		}

		// Print the certificate
		result, err = certinfo.CertificateText(cert)
		if err != nil {
			log.Printf("Failed to get cert info from URL error: %v", err)
			return err
		}
	}

	if c.IsSet(flagCert) {
		path := c.String(flagCert)
		cert, err = utils.CertFromFile(path)
		if err != nil {
			log.Printf("failed to get cert from file %q error: %v", path, err)
			return err
		}

		// Print the certificate
		result, err = certinfo.CertificateText(cert)
		if err != nil {
			log.Printf("Failed to get cert info from cert file %q error: %v", path, err)
			return err
		}
	}

	if c.IsSet(flagCSR) {
		path := c.String(flagCSR)
		csr, err = utils.CSRFromFile(path)
		if err != nil {
			log.Printf("failed to get CSR from file %q: %v", path, err)
			return err
		}

		// Print the certificate
		result, err = certinfo.CertificateRequestText(csr)
		if err != nil {
			log.Printf("Failed to get CSR info from cert error: %v", err)
			return err
		}
	}

	if err = os.WriteFile(outputFilePath, []byte(result), 0o600); err != nil {
		log.Printf("Failed to write PEM to file %s error: %v", outputFilePath, err)
		return err
	}
	return nil
}

func readX509FromDomain(uri string) (*x509.Certificate, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	// if no schema is used, parse with default scheme //
	// to get the host name
	if u.Host == "" {
		uri = "//" + uri
		u, err = url.Parse(uri)
		if err != nil {
			return nil, err
		}
	}

	uri = u.Host
	if u.Port() == "" {
		uri += ":443"
	}

	conn, err := tls.Dial("tcp", uri, &tls.Config{})
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	// Get certificates returned from the server
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		// Return the first certificate
		return certs[0], nil
	}

	return nil, fmt.Errorf("no certificate returned from %q", uri)
}
