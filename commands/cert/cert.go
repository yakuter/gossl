package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"io"
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
	flagIsCSR  = "isCSR"
)

func Command(reader io.Reader) *cli.Command {
	return &cli.Command{
		Name:        CmdCert,
		HelpName:    CmdCert,
		Action:      Action(reader),
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
			Value:       365,
			Required:    false,
		},
		&cli.Uint64Flag{
			Name:        flagSerial,
			Usage:       "Serial number to use in certificate",
			DefaultText: "123456",
			Value:       1,
			Required:    false,
		},
		&cli.BoolFlag{
			Name:     flagIsCA,
			Usage:    "Is Root Certificate Authority (CA) flag",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     flagIsCSR,
			Usage:    "Is Certificate Signing Request (CSR) flag",
			Required: false,
		},
	}
}

func Action(reader io.Reader) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		// Get privatekey from file
		privateKey, err := utils.PrivateKeyFromPEMFile(c.String(flagKey))
		if err != nil {
			log.Printf("Failed to get key from key file %s error: %v", c.String(flagKey), err)
			return err
		}

		// Set output
		output := os.Stdout
		outputFilePath := output.Name()
		if c.IsSet(flagOut) {
			outputFilePath = c.String(flagOut)
		}

		// Generate subject (pkix.Name) from answers
		subj, email, err := subject(reader)
		if err != nil {
			log.Printf("Failed to generate subject error: %v", err)
			return err
		}

		var outPEM []byte
		if c.Bool(flagIsCSR) {
			outPEM, err = generateCSR(subj, email, privateKey)
		} else {
			outPEM, err = generateCert(subj, c.Uint(flagDays), c.Uint64(flagSerial), c.Bool(flagIsCA), privateKey)
		}
		if err != nil {
			log.Printf("Failed to create cert error: %v", err)
			return err
		}

		// Write x509 certificate to file
		if err = os.WriteFile(outputFilePath, outPEM, 0o600); err != nil {
			log.Printf("Failed to write PEM to file %s error: %v", outputFilePath, err)
			return err
		}

		log.Printf("Certificate generated")
		return nil
	}
}

func subject(reader io.Reader) (pkix.Name, string, error) {
	// Prepare questions which are needed for subject
	questions := []string{
		"Common Name - SAN (eg, FQDN or IP)* []",
		"E-mail address* []",
		"Country Name (2 letter code) [AU]",
		"State or Province Name []",
		"Locality Name (eg, city) []",
		"Organization Name [eg, company]",
		"Organizational Unit Name (eg, section) []",
		"Street Addr []",
		"Postal Code []",
	}

	// Ask questions to user and get inputs as answers
	answers, err := utils.ReadInputs(questions, reader)
	if err != nil {
		log.Printf("failed to read inputs %v", err)
		return pkix.Name{}, "", err
	}

	san := answers[0]
	email := answers[1]

	if len(san) == 0 {
		err = errors.New("Common Name - SAN cannot be empty")
		log.Printf("%v", err)
		return pkix.Name{}, "", err
	}

	return pkix.Name{
		Country:            []string{answers[2]},
		Province:           []string{answers[3]},
		Locality:           []string{answers[4]},
		Organization:       []string{answers[5]},
		OrganizationalUnit: []string{answers[6]},
		StreetAddress:      []string{answers[7]},
		PostalCode:         []string{answers[8]},
	}, email, nil
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

func generateCert(subj pkix.Name, days uint, serial uint64, isCA bool, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Generate template (x509 certificate)
	t := template(subj, days, serial, isCA)

	// Create x509 certificate
	certx509, err := x509.CreateCertificate(rand.Reader, t, t, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to create certificate error: %v", err)
		return nil, err
	}

	// Return cert in PEM format
	return utils.CertToPEM(certx509), nil
}

func generateCSR(subj pkix.Name, email string, privateKey *rsa.PrivateKey) ([]byte, error) {
	if len(email) == 0 {
		err := errors.New("E-mail address cannot be empty")
		log.Printf("%v", err)
		return nil, err
	}

	oidEmailAddress := asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

	rawSubj := subj.ToRDNSequence()
	rawSubj = append(rawSubj, []pkix.AttributeTypeAndValue{
		{Type: oidEmailAddress, Value: email},
	})

	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		log.Printf("Failed to asn1.Marshal raw subject for CSR error: %v", err)
		return nil, err
	}

	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		EmailAddresses:     []string{email},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	// Create x509 certificate request (CSR)
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		log.Printf("Failed to create certificate request error: %v", err)
		return nil, err
	}

	// Return CSR in PEM format
	return utils.CSRToPEM(csrBytes), nil
}
