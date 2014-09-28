package ttunnel

import (
	"fmt"
	"os"
	"path/filepath"
)

// Helper function for getting the appropriate path for a client's
// configuration file by name.
func ClientConfigPath(name string) string {
	return filepath.Join(TunnelDir, name+".json")
}

// Helper function for getting the removed client's configuration file
// by name.
func ClientRemovedPath(name string) string {
	return filepath.Join(RemovedDir, name+".json")
}

// ClientConfig represents a client's configuration file. A single
// configuration file will allow a connection with a single host on a
// single port. The local port can be chosen arbitrarily.
type ClientConfig struct {
	Host  string // The host: <address>:<port>
	Port  int32  // The local port to listen on.
	Token []byte // The token.
}

// ReadClientConfig reads the named configuration from
// $(HOME)/.ttunnel/tunnels/<name>.
func ReadClientConfig(name string) (cc ClientConfig, err error) {
	UnmarshalFrom(ClientConfigPath(name), & cc)
	return
}

// WriteClientConfig writes the given configuration file to
// $(HOME)/.ttunnel/tunnels/<name>.
func WriteClientConfig(name string, cc ClientConfig) (err error) {
	err = MarshalTo(ClientConfigPath(name), cc)
	return
}

// RemoveClientConfig removes a file in $(HOME)/.ttunnel/tunnels/ for
// the named client.
func RemoveClientConfig(name string) error {
	fromPath := ClientConfigPath(name)
	toPath := ClientRemovedPath(name)

	// Skip if file doesn't exist.
	if !FileExists(fromPath) {
		return nil
	}

	// Try not to clobber.
	if FileExists(toPath) {
		return fmt.Errorf("Not overwriting file: %v", toPath)
	}

	return os.Rename(fromPath, toPath)
}
