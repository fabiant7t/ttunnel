package ttunnel

import (
	"crypto/rand"
	"encoding/hex"
)

func randomSecret() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// AddClient creates the named client tunnel configuration.
func AddClient(
	name, hostAddr, connectAddr string, localPort int32) (err error) {

	// Create the client's tunnel configuration.
	tc := TunnelConfig{}
	tc.Host = hostAddr
	tc.Port = localPort

	tc.Pwd, err = randomSecret()
	if err != nil {
		return err
	}

	_, tc.CaCert, err = loadKeyAndCert()
	if err != nil {
		return err
	}

	// Create the server's configuration file.
	cc := ClientConfig{}
	cc.ConnectAddr = connectAddr
	cc.Secret = tc.Pwd

	// Write the config files.
	if err = tc.Save(clientTunnelPath(name)); err != nil {
		return err
	}

	return cc.Save(clientsPath(name))
}
