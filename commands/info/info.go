package info

import (
	"errors"
	"fmt"
	"log"

	"github.com/yakuter/gossl/pkg/utils"

	"github.com/grantae/certinfo"
	"github.com/urfave/cli/v2"
)

const (
	CmdInfo = "info"
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
	return []cli.Flag{}
}

func Action(c *cli.Context) error {
	if c.Args().Len() == 0 {
		err := errors.New("cert file argument is not found")
		log.Printf("%v", err)
		return err
	}

	// Get certificate from file
	certFilePath := c.Args().First()
	cert, err := utils.CertFromFile(certFilePath)
	if err != nil {
		log.Printf("Failed to get cert from file %s CAs error: %v", certFilePath, err)
		return err
	}

	// Print the certificate
	result, err := certinfo.CertificateText(cert)
	if err != nil {
		log.Printf("Failed to get cert info from cert error: %v", err)
		return err
	}

	fmt.Println(result)
	return nil
}
