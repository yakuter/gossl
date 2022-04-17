package info

import (
	"errors"
	"log"
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
	}
}

func Action(c *cli.Context) error {
	if c.Args().Len() == 0 {
		err := errors.New("cert file argument is not found")
		log.Printf("%v", err)
		return err
	}

	// Set output
	output := os.Stdout
	outputFilePath := output.Name()
	if c.IsSet(flagOut) {
		outputFilePath = c.String(flagOut)
	}

	// Get certificate from file
	certFilePath := c.Args().First()
	cert, err := utils.CertFromFile(certFilePath)
	if err != nil {
		log.Printf("Failed to get cert from file %s error: %v", certFilePath, err)
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
