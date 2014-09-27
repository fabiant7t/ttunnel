package ttunnel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	buf, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return
	}

	if err = json.Unmarshal(buf, &sc); err != nil {
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
	if FileExists(ConfigPath) {
		err = fmt.Errorf("Won't overwrite file: %v", ConfigPath)
		return
	}

	if len(sc.EncKey) != 32 {
		err = fmt.Errorf("Key must be 32 bytes, not %v bytes.", len(sc.EncKey))
		return
	}

	buf, err := json.Marshal(sc)
	if err != nil {
		return
	}

	var out bytes.Buffer
	if err = json.Indent(&out, buf, "", "\t"); err != nil {
		return
	}

	err = ioutil.WriteFile(ConfigPath, out.Bytes(), 0600)
	return
}
