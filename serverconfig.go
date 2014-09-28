package ttunnel

import (
	"fmt"
)

// ServerConfig represents a server's configuration file.
type ServerConfig struct {
	ListenAddr  string // The address to listen on. For example :1044.
	ConnectAddr string // The address client's would connect to. blah.com:1044
	EncKey      []byte // The encryption key for signing and encryption.
}

// ReadServerConfig reads the server configuration file from the
// appropriate location.
func ReadServerConfig() (sc ServerConfig, err error) {
	if err = UnmarshalFrom(ConfigPath, &sc); err != nil {
		return
	}

	if len(sc.EncKey) != 32 {
		err = fmt.Errorf("Key must be 32 bytes, not %v bytes.", len(sc.EncKey))
	}

	return
}

// WriteServerConfig writes the given server configuration file to the
// appropriate location.
func WriteServerConfig(sc ServerConfig) (err error) {
	if len(sc.EncKey) != 32 {
		err = fmt.Errorf("Key must be 32 bytes, not %v bytes.", len(sc.EncKey))
		return
	}
	
	err = MarshalTo(ConfigPath, sc)
	return
}
