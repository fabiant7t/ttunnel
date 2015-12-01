package ttunnel

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"path/filepath"
)

func RunClient() error {
	matches, err := filepath.Glob(filepath.Join(tunnelDir, "*.json"))
	if err != nil {
		return err
	}

	for _, path := range matches {
		_, name := filepath.Split(path)
		name = name[:len(name)-5]
		go runClient(name)
	}

	// Sleep forever.
	select {}
}

// RunClient runs a client using the named configuration.
func runClient(name string) {
	log.Printf("Loading tunnel: %v\n", name)

	// Load the tunnel configuration file.
	tc := TunnelConfig{}
	if err := tc.Load(tunnelPath(name)); err != nil {
		log.Printf("Error when loading config file: %v\n", err)
		return
	}

	// Create a certificate pool for the client.
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(tc.CaCert) {
		log.Printf("Failed to append certificate to pool.")
		return
	}

	// Create the client configuration.
	config := tls.Config{RootCAs: certPool}

	// Accept connections on the local port.
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", tc.Port))
	if err != nil {
		log.Printf("Failed to list on TCP port %v: %v\n", tc.Port, err)
		return
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
}

func clientHandler(
	lConn net.Conn, host, name string, pwd string, config tls.Config) {

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
	if err = writeBytes(rConn, []byte(pwd)); err != nil {
		log.Printf("Error sending secret:\n    %v\n", err)
		goto closeConns
	}

	// Forward traffic.
	copyDuplex(lConn, rConn)
	return

closeConns:
	rConn.Close()
	lConn.Close()
}
