package cert

import (
	"bufio"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/urfave/cli/v2"
)

const (
	CmdCert = "cert"

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
			Name:        flagOut,
			Usage:       "Output file name (optional)",
			DefaultText: "eg, ./cert.pem",
			Required:    false,
		},
		&cli.UintFlag{
			Name:        flagDays,
			Usage:       "Number of days a certificate is valid for",
			DefaultText: "365",
			Required:    true,
		},
		&cli.Uint64Flag{
			Name:        flagSerial,
			Usage:       "Serial number to use in certificate",
			DefaultText: "123456",
			Required:    false,
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
		err = errors.New("Common Name - SAN cannot be empty")
		log.Printf("%v", err)
		return err
	}

	// Generate subject (pkix.Name) from answers
	p := subject(answers)

	// Generate template (x509 certificate)
	t := template(p, c.Uint(flagDays), c.Uint64(flagSerial), c.Bool(flagIsCA))

	// Get privatekey from file
	privateKey, err := utils.PrivateKeyFromPEMFile(c.String(flagKey))
	if err != nil {
		log.Printf("Failed to get key from key file %s error: %v", c.String(flagKey), err)
		return err
	}

	// Create x509 certificate
	certx509, err := x509.CreateCertificate(rand.Reader, t, t, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to create certificate error: %v", err)
		return err
	}

	// Encode x509 certificate to PEM format
	certBytes := utils.CertToPEM(certx509)

	// Set output
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	// Write x509 certificate to file
	if err = os.WriteFile(outputFilePath, certBytes, 0o600); err != nil {
		log.Printf("Failed to write Public Key to file %s error: %v", outputFilePath, err)
		return err
	}

	log.Printf("Certificate generated")
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
