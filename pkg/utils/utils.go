package utils

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// GeneratePublicKey take an rsa.PublicKey and return bytes
// suitable for writing to .pub file in the format "ssh-rsa ..."
func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(publicRsaKey), nil
}

// GeneratePrivateKey creates an RSA Private Key with provided bit size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	if err = privateKey.Validate(); err != nil {
		return nil, err
	}

	return privateKey, nil
}

// PrivateKeyToPEM encodes Private Key from RSA to PEM format
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// CertToPEM encodes Certificate to PEM format
func CertToPEM(cert []byte) []byte {
	// pem.Block
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	// Cert in PEM format
	certPEM := pem.EncodeToMemory(&block)

	return certPEM
}

func PrivateKeyFromPEMFile(keyFilePath string) (*rsa.PrivateKey, error) {
	keyFileContent, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}

	keyBlock, _ := pem.Decode(keyFileContent)
	if keyBlock == nil {
		return nil, fmt.Errorf("invalid key file %s", keyFilePath)
	}

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func ReadInputs(questions []string) ([]string, error) {
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
