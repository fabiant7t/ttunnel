package ttunnel

import (
	"github.com/johnnylee/goutil/jsonutil"
)

// A ClientConfig represents the configuration file for a tunnel
// client, stored on the server.
type ClientConfig struct {
	ConnectAddr string // Standard format: <address>:<port>
	Secret      string // Shared secret key.
}

// Load loads client configuration from the given file.
func (cc *ClientConfig) Load(path string) error {
	return jsonutil.Load(cc, path)
}

// Save stores the client configuration into the given file.
func (cc ClientConfig) Save(path string) error {
	return jsonutil.Store(cc, path)
}
