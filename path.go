package ttunnel

import (
	"github.com/johnnylee/goutil/fileutil"
)

var configDir string       // The base configuration directory.
var clientsDir string      // The server's clients directory.
var clientTunnelDir string // The server's client file directory.
var tunnelDir string       // The client's tunnel directory.
var keyPath string         // The path to the server's private key.
var certPath string        // The path to the server's certificate.

// init will initialize the global paths in this file.
func init() {
	// Initialize directory paths.
	configDir = fileutil.ExpandPath("~/.ttunnel")
	clientsDir = fileutil.ExpandPath(configDir, "clients")
	clientTunnelDir = fileutil.ExpandPath(configDir, "client-tunnels")
	tunnelDir = fileutil.ExpandPath(configDir, "tunnels")

	// Initialize file paths.
	certPath = fileutil.ExpandPath(configDir, "server.crt")
	keyPath = fileutil.ExpandPath(configDir, "server.key")
}

func clientsPath(name string) string {
	return fileutil.ExpandPath(clientsDir, name+".json")
}

func clientTunnelPath(name string) string {
	return fileutil.ExpandPath(clientTunnelDir, name+".json")
}

func tunnelPath(name string) string {
	return fileutil.ExpandPath(tunnelDir, name+".json")
}
