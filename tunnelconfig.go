package ttunnel

// A TunnelConfig represents the configuration of a tunnel from the
// client's perspective.
type TunnelConfig struct {
	Host   string // The host address: <address>:<port>
	Port   int32  // The local port to listen on.
	Secret string // The shared secret.
	CaCert []byte // The CA certificate.
}

// Load loads the tunnel configuration from the given file.
func (tc *TunnelConfig) Load(path string) error {
	return UnmarshalFromFile(path, tc)
}

// Save stores the tunnel configuration into the given file.
func (tc TunnelConfig) Save(path string) error {
	return MarshalToFile(path, &tc)
}
