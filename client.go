package ttunnel

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
)

// RunClient runs a client using the named configuration.
func RunClient(name string) error {
	// Load the tunnel configuration file.
	tc := TunnelConfig{}
	if err := tc.Load(tunnelPath(name)); err != nil {
		return err
	}

	// Create a certificate pool for the client.
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(tc.CaCert) {
		return fmt.Errorf("Error appending certificate to pool.")
	}

	// Create the client configuration.
	config := tls.Config{RootCAs: certPool}

	// Accept connections on the local port.
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", tc.Port))
	if err != nil {
		return err
	}

	// Accept connections forever.
	for {
		lConn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection:\n    %v\n", err)
			continue
		}

		go clientHandler(lConn, tc.Host, name, tc.Pwd, config)
	}

	return nil
}

func clientHandler(
	lConn net.Conn, host, name string, pwd []byte, config tls.Config) {

	// Dial remote server.
	rConn, err := tls.Dial("tcp", host, &config)
	if err != nil {
		log.Printf("Error dialing server:\n    %v\n", err)
		lConn.Close()
		return
	}

	// Send name.
	// Send the name.
	if err = writeBytes(rConn, []byte(name)); err != nil {
		log.Printf("Error sending name:\n    %v\n", err)
		goto closeConns
	}

	// Send the password.
	if err = writeBytes(rConn, pwd); err != nil {
		log.Printf("Error sending password:\n    %v\n", err)
		goto closeConns
	}

	// Forward traffic.
	copyDuplex(lConn, rConn)
	return

closeConns:
	rConn.Close()
	lConn.Close()
}
