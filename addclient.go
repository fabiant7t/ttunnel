package ttunnel

import (
	"code.google.com/p/go.crypto/bcrypt"
	"crypto/rand"
)

// AddClient creates the named client tunnel configuration.
func AddClient(
	name, hostAddr, connectAddr string, localPort int32) (err error) {

	// Create the client's tunnel configuration.
	tc := TunnelConfig{}
	tc.Host = hostAddr
	tc.Pwd = make([]byte, 48)
	tc.Port = localPort

	if _, err = rand.Read(tc.Pwd); err != nil {
		return
	}

	_, tc.CaCert, err = loadKeyAndCert()
	if err != nil {
		return
	}

	// Create the server's configuration file.
	cc := ClientConfig{}
	cc.ConnectAddr = connectAddr
	cc.PwdHash, err = bcrypt.GenerateFromPassword(tc.Pwd, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	// Write the config files.
	if err = tc.Save(clientTunnelPath(name)); err != nil {
		return
	}

	err = cc.Save(clientsPath(name))

	return
}
