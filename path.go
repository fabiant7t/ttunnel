package ttunnel

import (
	"os"
	"os/user"
	"path/filepath"
)

var configDir string       // The base configuration directory.
var clientsDir string      // The server's clients directory.
var clientTunnelDir string // The server's client file directory.
var tunnelDir string       // The client's tunnel directory.
var keyPath string         // The path to the server's private key.
var certPath string        // The path to the server's certificate.

// init will initialize the global paths in this file.
func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Initialize directory paths.
	configDir = filepath.Join(usr.HomeDir, ".ttunnel")
	clientsDir = filepath.Join(configDir, "clients")
	clientTunnelDir = filepath.Join(configDir, "client-tunnels")
	tunnelDir = filepath.Join(configDir, "tunnels")

	// Initialize file paths.
	certPath = filepath.Join(configDir, "server.crt")
	keyPath = filepath.Join(configDir, "server.key")
}

// fileExists checks if a path exists. Arguments are joined with
// filepath.Join to construct the full path.
func fileExists(elem ...string) bool {
	path := filepath.Join(elem...)
	_, err := os.Stat(path)
	return err == nil
}

func clientsPath(name string) string {
	return filepath.Join(clientsDir, name+".json")
}

func clientTunnelPath(name string) string {
	return filepath.Join(clientTunnelDir, name+".json")
}

func tunnelPath(name string) string {
	return filepath.Join(tunnelDir, name+".json")
}
