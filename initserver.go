package ttunnel

import (
	"os"
)

// InitializeServer creates the directories and files necessary to
// operate the server, including the key and certificate files.
func InitServer(host string, rsaBits int) error {
	err := os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(clientsDir, 0700); err != nil {
		return err
	}

	if err := os.MkdirAll(clientTunnelDir, 0700); err != nil {
		return err
	}

	err = generateKeyAndCert(host, true, rsaBits)
	if err != nil {
		return err
	}

	return nil
}
