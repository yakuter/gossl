package utils

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

// GenerateSSHPublicKey take an rsa.PublicKey and return bytes
// suitable for writing to .pub file in the format "ssh-rsa ..."
func GenerateSSHPublicKey(rsaPubKey *rsa.PublicKey) ([]byte, error) {
	sshPubKey, err := ssh.NewPublicKey(rsaPubKey)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(sshPubKey), nil
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

// PrivateKeyToPEM encodes Private Key to PEM format
func PrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}

	return pem.EncodeToMemory(&privBlock)
}

// PublicKeyToPEM encodes Public Key to PEM format
func PublicKeyToPEM(publicKey *rsa.PublicKey) []byte {
	// Get ASN.1 DER format
	pubDER := x509.MarshalPKCS1PublicKey(publicKey)

	// pem.Block
	pubBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	}

	return pem.EncodeToMemory(&pubBlock)
}

// CertToPEM encodes Certificate to PEM format
func CertToPEM(cert []byte) []byte {
	// pem.Block
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	return pem.EncodeToMemory(&block)
}

// CSRToPEM encodes CSR to PEM format
func CSRToPEM(cert []byte) []byte {
	// pem.Block
	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: cert,
	}

	return pem.EncodeToMemory(&block)
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

func CertFromFile(certFilePath string) (*x509.Certificate, error) {
	// Read cert file
	certFileBytes, err := os.ReadFile(certFilePath)
	if err != nil {
		log.Printf("Failed to read cert file %q error: %v", certFilePath, err)
		return nil, err
	}

	// Decode PEM encoded cert file
	block, _ := pem.Decode(certFileBytes)
	if block == nil {
		err = errors.New("block is nil")
		log.Printf("Failed to decode PEM encoded cert file %q error: %v", certFilePath, err)
		return nil, err
	}

	// Parse x509 certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Printf("Failed to parse x509 certificate from cert file %q error: %v", certFilePath, err)
		return nil, err
	}

	return cert, nil
}

func ReadInputs(questions []string, reader io.Reader) ([]string, error) {
	answers := make([]string, len(questions))
	scanner := bufio.NewScanner(reader)
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
