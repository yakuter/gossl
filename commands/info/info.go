package info

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/grantae/certinfo"
	"github.com/urfave/cli/v2"
)

const (
	CmdInfo = "info"

	flagOut = "out"
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
	}
}

func Action(c *cli.Context) error {
	if c.Args().Len() == 0 {
		err := errors.New("cert file or URL argument is not found")
		log.Printf("%v", err)
		return err
	}

	// Set output
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	// Get certificate from file or URL
	certPath := c.Args().First()
	cert, err := readX509FromFileOrURL(certPath)
	if err != nil {
		return err
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

func readX509FromFileOrURL(path string) (*x509.Certificate, error) {
	// Check if given URL is a valid URI
	if u, err := url.ParseRequestURI(path); err == nil {
		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		// We don't need body but it must be read to EOF
		// before closing.
		io.Copy(ioutil.Discard, resp.Body)

		// Get certificate returned from the server
		if resp.TLS != nil {
			certs := resp.TLS.PeerCertificates
			if len(certs) > 0 {
				return certs[0], nil
			}
		}

		return nil, fmt.Errorf("no certificate returned from %q", u.String())
	}

	// Read from file
	cert, err := utils.CertFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get cert from file %q: %v", path, err)
	}

	return cert, nil
}
