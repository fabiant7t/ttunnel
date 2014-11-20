package ttunnel

import (
	"code.google.com/p/go.crypto/bcrypt"
)

// A ClientConfig represents the configuration file for a tunnel
// client, stored on the server.
type ClientConfig struct {
	ConnectAddr string // Standard format: <address>:<port>
	PwdHash     []byte // Bcrypt-encrypted password hash.
}

// Load loads client configuration from the given file.
func (cc *ClientConfig) Load(path string) error {
	return UnmarshalFromFile(path, cc)
}

// Save stores the client configuration into the given file.
func (cc ClientConfig) Save(path string) error {
	return MarshalToFile(path, &cc)
}

// PwdMatches returns true if the give password matches the stored
// password.
func (cc ClientConfig) PwdMatches(pwd []byte) bool {
	return bcrypt.CompareHashAndPassword(cc.PwdHash, pwd) == nil
}
