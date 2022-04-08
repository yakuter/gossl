package cert

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdCert = "cert"
)

const (
	flagKey    = "key"
	flagOut    = "out"
	flagDays   = "days"
	flagSerial = "serial"
	flagIsCA   = "isCA"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:        CmdCert,
		HelpName:    CmdCert,
		Action:      Action,
		ArgsUsage:   `[cert file path]`,
		Usage:       `generates x509 certificate.`,
		Description: `Generates x509 certificate with provided template information.`,
		Flags:       Flags(),
	}
}

func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagKey,
			Usage:    "private key (required)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagOut,
			Usage:    "Output file name (optional)",
			Required: false,
		},
		&cli.UintFlag{
			Name:     flagDays,
			Usage:    "Number of days a certificate is valid for",
			Required: true,
		},
		&cli.Uint64Flag{
			Name:     flagSerial,
			Usage:    "Serial number to use in certificate",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     flagIsCA,
			Usage:    "Number of days the certificate is valid for",
			Required: false,
		},
	}
}

func Action(c *cli.Context) error {
	questions := []string{
		"Common Name - SAN (eg, FQDN or IP)* []",
		"Country Name (2 letter code) [AU]",
		"State or Province Name []",
		"Locality Name (eg, city) []",
		"Organization Name [eg, company]",
		"Organizational Unit Name (eg, section) []",
		"Street Addr []",
		"Postal Code []",
	}

	// Ask questions to user and get answers
	answers, err := readInputs(questions)
	if err != nil {
		log.Printf("failed to read inputs %v", err)
		return err
	}

	if len(answers[0]) == 0 {
		err := errors.New("Common Name - SAN cannot be empty")
		log.Printf("%v", err)
		return err
	}

	// Generate subject (pkix.Name) from answers
	p := subject(answers)

	// Generate template (x509 certificate)
	t := template(p, c.Uint(flagDays), c.Uint64(flagSerial), c.Bool(flagIsCA))

	// Get privatekey from file
	privateKey, err := key(c.String(flagKey))
	if err != nil {
		log.Printf("Failed to get key from key file %s error: %v", c.String(flagKey), err)
		return err
	}

	// Create x509 certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, t, t, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to create certificate error: %v", err)
		return err
	}

	// Encode certificate as PEM
	certPEM := bytes.NewBuffer(nil)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		log.Printf("Failed to encode cert as pem error: %v", err)
		return err
	}

	// Set output
	var output *os.File = os.Stdout

	// If output file is provided, then create it and set as output
	if c.IsSet(flagOut) {
		outputFilePath := c.String(flagOut)
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			log.Printf("failed to create output file: %q error: %v", outputFilePath, err)
			return err
		}

		defer func() {
			if err := outputFile.Close(); err != nil {
				log.Printf("failed to close output file: %q error: %v", outputFilePath, err)
			}
		}()

		output = outputFile
	}

	// Write PEM encoded x509 certificate to output
	_, err = output.WriteString(certPEM.String())
	if err != nil {
		log.Printf("failed to write output error: %v", err)
		return err
	}

	log.Printf("Certificate generated successfully")
	return nil
}

func readInputs(questions []string) ([]string, error) {
	answers := make([]string, len(questions))
	scanner := bufio.NewScanner(os.Stdin)
	for i := range questions {
		fmt.Printf("%s: ", questions[i])
		if ok := scanner.Scan(); !ok {
			return nil, errors.New("failed to scan")
		}
		answers[i] = scanner.Text()
	}

	if scanner.Err() != nil {
		log.Printf("Scanner error: %v", scanner.Err())
	}

	return answers, nil
}

func subject(answers []string) pkix.Name {
	return pkix.Name{
		Country:            []string{answers[1]},
		Province:           []string{answers[2]},
		Locality:           []string{answers[3]},
		Organization:       []string{answers[4]},
		OrganizationalUnit: []string{answers[5]},
		StreetAddress:      []string{answers[6]},
		PostalCode:         []string{answers[7]},
	}
}

func template(subject pkix.Name, days uint, serial uint64, isCA bool) *x509.Certificate {
	t := &x509.Certificate{
		SerialNumber: big.NewInt(int64(serial)),
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, int(days)),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}

	if isCA {
		t.IsCA = true
		t.BasicConstraintsValid = true
	}
	return t
}

func key(keyFilePath string) (*rsa.PrivateKey, error) {
	keyFileContent, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyFileContent)
	if keyBlock == nil {
		err := fmt.Errorf("invalid key file %s", keyFilePath)
		log.Printf("%v", err)
		return nil, err
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		log.Printf("Failed to parse key file as x509 PKCS1 form error: %v", err)
		return nil, err
	}

	return key, nil
}
