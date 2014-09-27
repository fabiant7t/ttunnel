package ttunnel

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var ConfigDir string  // Configuration directory.
var ConfigPath string // The server configuration file path.
var TunnelDir string  // Directory containing tunnels.
var RemovedDir string // Removed client directory.
var RootCAPath string // Path to rootCA.crt file.
var CertPath string   // Path to X.509 certificate file.
var KeyPath string    // Path to X.509 key file.

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Initialize directory paths
	ConfigDir = filepath.Join(usr.HomeDir, ".ttunnel")
	ConfigPath = filepath.Join(ConfigDir, "server.json")
	TunnelDir = filepath.Join(ConfigDir, "tunnels")
	RemovedDir = filepath.Join(ConfigDir, "removed")

	// Create directories.
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(TunnelDir, 0700); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(RemovedDir, 0700); err != nil {
		panic(err)
	}

	// Initialize file paths.
	RootCAPath = filepath.Join(ConfigDir, "rootCA.crt")
	CertPath = filepath.Join(ConfigDir, "server.crt")
	KeyPath = filepath.Join(ConfigDir, "server.key")
}

// fileExists checks if a path exists. Arguments are joined with
// filepath.Join to construct the full path.
func FileExists(elem ...string) bool {
	path := filepath.Join(elem...)
	_, err := os.Stat(path)
	return err == nil
}

// RandomBytes return cryptographically secure random bytes.
func RandomBytes(length int) (key []byte, err error) {
	key = make([]byte, length)
	_, err = io.ReadFull(rand.Reader, key)
	return
}

// loadRootCA loads the custom root certificate authority if the
// $(HOME)/.ttunnel/rootCA.crt file exists.
func loadRootCA(config *tls.Config) error {

	// If the file doesn't exist, we'll use the system's default.
	if !FileExists(RootCAPath) {
		return nil
	}

	// Load the custom rootCA.
	config.RootCAs = x509.NewCertPool()

	rootCert, err := ioutil.ReadFile(RootCAPath)
	if err != nil {
		return err
	}

	if !config.RootCAs.AppendCertsFromPEM(rootCert) {
		return fmt.Errorf("Failed to load root CA certificate.")
	}

	return nil
}
