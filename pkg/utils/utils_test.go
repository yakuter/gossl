package utils_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yakuter/gossl/pkg/utils"
)

func TestGenerateSSHPublicKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pubBytes, err := utils.GenerateSSHPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)
	require.NotNil(t, pubBytes)
	require.True(t, strings.HasPrefix(string(pubBytes), "ssh-rsa"))
}

func TestGeneratePrivateKey(t *testing.T) {
	key, err := utils.GeneratePrivateKey(2048)
	require.NoError(t, err)
	require.NotNil(t, key)
}

func TestPrivateKeyToPEM(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pem := utils.PrivateKeyToPEM(privateKey)
	require.NotNil(t, pem)
	require.True(t, strings.Contains(string(pem), "RSA PRIVATE KEY"))
}

func TestPublicKeyToPEM(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pem := utils.PublicKeyToPEM(&privateKey.PublicKey)
	require.NotNil(t, pem)
	require.Contains(t, string(pem), "RSA PUBLIC KEY")
}

func TestCertToPEM(t *testing.T) {
	pem := utils.CertToPEM([]byte("testCert"))
	require.NotNil(t, pem)
	require.Contains(t, string(pem), "CERTIFICATE")
}

func TestCSRToPEM(t *testing.T) {
	pem := utils.CSRToPEM([]byte("testCert"))
	require.NotNil(t, pem)
	require.Contains(t, string(pem), "CERTIFICATE REQUEST")
}

func TestPrivateKeyFromPEMFile(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	file, err := os.CreateTemp(t.TempDir(), "private-*.key")
	require.NoError(t, err)

	// Fail
	keyFromFile, err := utils.PrivateKeyFromPEMFile("wrong-file")
	require.Error(t, err)

	// Fail
	keyFromFile, err = utils.PrivateKeyFromPEMFile(file.Name())
	require.Error(t, err)

	// Write test pem to file
	_, err = file.Write(utils.PrivateKeyToPEM(privateKey))
	require.NoError(t, err)
	require.NoError(t, file.Close())

	// Success
	keyFromFile, err = utils.PrivateKeyFromPEMFile(file.Name())
	require.NoError(t, err)
	require.True(t, reflect.DeepEqual(privateKey, keyFromFile))
}

func TestReadInputs(t *testing.T) {
	q := []string{"Question1", "Question2"}

	var stdin bytes.Buffer

	// Success
	stdin.Write([]byte("localhost\njohn@doe.com\n"))
	answers, err := utils.ReadInputs(q, &stdin)
	require.NoError(t, err)
	require.Equal(t, answers, []string{"localhost", "john@doe.com"})

	// Success
	stdin.Write([]byte("localhost\njohn@doe.com\njane@doe.com\n"))
	answers, err = utils.ReadInputs(q, &stdin)
	require.NoError(t, err)
	require.Equal(t, answers, []string{"localhost", "john@doe.com"})

	// Fail
	stdin.Write([]byte("localhost\n"))
	answers, err = utils.ReadInputs(q, &stdin)
	require.Error(t, err)
}

func TestCertFromFile(t *testing.T) {
	// Load a valid certificate
	const validCertPath = "../../testdata/server-cert.pem"
	_, err := utils.CertFromFile(validCertPath)
	require.NoError(t, err)

	_, err = utils.CertFromFile("wrong-path")
	require.Error(t, err)

	// Fails decoding PEM
	invalidFile, err := os.CreateTemp(t.TempDir(), "invalid-*.pem")
	_, err = utils.CertFromFile(invalidFile.Name())
	require.Error(t, err)
}
