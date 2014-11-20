package ttunnel

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

// Function that generates certificates. Mostly coppied from
// http://golang.org/src/pkg/crypto/tls/generate_cert.go.
func generateKeyAndCert(host string, isCA bool, rsaBits int) error {

	// Make sure we have at least a 2048 bit key.
	if rsaBits < 2048 {
		return fmt.Errorf("Invalid number of bits for RSA: %v", rsaBits)
	}

	// Generate the private key.
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return err
	}

	// Generate a random serial number.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	// Create a certificate template.
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{"None"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * 365 * 100 * time.Hour),
		KeyUsage: (x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Set the host name.
	template.DNSNames = append(template.DNSNames, host)

	// Set certificate authority flags.
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	// Create the certificate.
	derBytes, err := x509.CreateCertificate(
		rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Encode and save the key.
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	err = pem.Encode(
		keyOut,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		})
	if err != nil {
		return err
	}

	// Encode and save the certificate.
	certOut, err := os.OpenFile(certPath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer certOut.Close()

	err = pem.Encode(
		certOut,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		})
	if err != nil {
		return err
	}

	return nil
}

// loadKeyAndCert reads the key and certificate files into memory and
// returns them in that order.
func loadKeyAndCert() ([]byte, []byte, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, nil, err
	}

	return key, cert, err
}
