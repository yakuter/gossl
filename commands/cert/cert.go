package cert

import (
	"bufio"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

// Remote commands
const (
	CmdCert = "cert"
)

const (
	flagKey  = "key"
	flagOut  = "out"
	flagDays = "days"
	flagIsCA = "isCA"
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
		&cli.BoolFlag{
			Name:     flagIsCA,
			Usage:    "Number of days a certificate is valid for",
			Required: false,
		},
	}
}

func Action(c *cli.Context) error {
	log.Printf("Certificate command args: %q\n", c.Args().Slice())

	answers, err := readInputs()
	if err != nil {
		return errors.Wrapf(err, "failed to read inputs")
	}

	result := setPkixName(answers)

	log.Printf("%q", result)

	log.Printf("Certificate generated successfully")
	return nil
}

func readInputs() ([]string, error) {
	answers := make([]string, 8)
	questions := []string{
		"Common Name - SAN (eg, server FQDN or IP) []: ",
		"Country Name (2 letter code) [AU]: ",
		"State or Province Name [Some-State]: ",
		"Locality Name (eg, city) []: ",
		"Organization Name [Internet Widgits Pty Ltd]: ",
		"Organizational Unit Name (eg, section) []: ",
		"Street Addr []: ",
		"Postal Code []: ",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for i := range questions {
		fmt.Print(questions[i])
		scanner.Scan()
		text := scanner.Text()
		if i == 0 && len(text) == 0 {
			err := errors.New("Common Name - SAN cannot be empty")
			return nil, errors.Wrapf(err, "error ")
		}
		answers[i] = text
	}

	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}

	return answers, nil
}

func setPkixName(answers []string) pkix.Name {
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
