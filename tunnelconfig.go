package ttunnel

import (
	"github.com/johnnylee/goutil/jsonutil"
)

// A TunnelConfig represents the configuration of a tunnel from the
// client's perspective.
type TunnelConfig struct {
	Host   string // The host address: <address>:<port>
	Port   int32  // The local port to listen on.
	Pwd    string // The shared secret.
	CaCert []byte // The CA certificate.
}

// Load loads the tunnel configuration from the given file.
func (tc *TunnelConfig) Load(path string) error {
	return jsonutil.Load(tc, path)
}

// Save stores the tunnel configuration into the given file.
func (tc TunnelConfig) Save(path string) error {
	return jsonutil.Store(tc, path)
}
