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

	flagOut = "out"
	flagURL = "url"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdInfo,
		HelpName:    CmdInfo,
		Action:      Action,
		ArgsUsage:   `[cert file path or URL ]`,
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
	}
}

func Action(c *cli.Context) error {
	// Set output
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	var (
		err  error
		cert *x509.Certificate
	)

	if c.IsSet(flagURL) {
		u := c.String(flagURL)
		cert, err = readX509FromDomain(u)
		if err != nil {
			log.Printf("failed to get cert details from %q: %v", u, err)
			return err
		}
	} else {
		if c.Args().Len() == 0 {
			err := errors.New("cert file argument is not found")
			log.Printf("%v", err)
			return err
		}

		// Get certificate from file
		path := c.Args().First()
		if err != nil {
			return err
		}

		// Read from file
		cert, err = utils.CertFromFile(path)
		if err != nil {
			log.Printf("failed to get cert from file %q: %v", path, err)
			return err
		}
	}

	// Print the certificate
	result, err := certinfo.CertificateText(cert)
	if err != nil {
		log.Printf("Failed to get cert info from cert error: %v", err)
		return err
	}

	// Write x509 certificate to file
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

	return nil, fmt.Errorf("no certificate returned from %q", u.String())
}
