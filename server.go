package ttunnel

import (
	"crypto/tls"
	"encoding/binary"
	"log"
	"net"
	"time"
)

// loadTlsConfig loads the TLS configuration. If
// $(HOME)/.ttunnel/rootCA.crt exists, then it will be loaded as the
// root certificate. If not, the defaults are used.
func loadTlsConfig() (config *tls.Config, err error) {
	config = new(tls.Config)

	// Load our server's cert.
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(CertPath, KeyPath)

	return
}

// handleConnection
func handleServerConnection(rConn net.Conn, th *TokenHandler) {
	// Read the client's encrypted token.
	var length uint16
	err := binary.Read(rConn, binary.LittleEndian, &length)
	if err != nil {
		log.Printf("Error reading token length:\n    %v\n", err)
		rConn.Close()
		return
	}

	encToken := make([]byte, length)
	_, err = rConn.Read(encToken)
	if err != nil {
		log.Printf("Error reading encrypted token:\n    %v\n", err)
		rConn.Close()
		return
	}

	// Decode the client's encrypted token.
	token, err := th.Decode(encToken)
	if err != nil {
		log.Printf("Error decoding token:\n    %v\n", err)
		rConn.Close()
		return
	}

	// Verify the token.
	if err = th.Verify(token); err != nil {
		log.Printf("Failed to verity token:\n    %v\n", err)
		rConn.Close()
		return
	}

	// Check the expiration date.
	dt := token.Expires - time.Now().Unix()
	days := dt / (3600 * 24)
	log.Printf("Login from %v. Token expires in %v days.\n", token.Name, days)

	// Create forwarded connection.
	fConn, err := net.Dial("tcp", token.ConnectAddr)
	if err != nil {
		log.Printf("Failed to make forwarding connection:\n    %v\n", err)
		rConn.Close()
		return
	}

	// Forward connections.
	forwardConnections(rConn, fConn)
}

// RunServer will run in the foreground forever if there are no errors
// initializing the server.
func RunServer() error {

	// Load the server's configuration.
	sc, err := ReadServerConfig()
	if err != nil {
		return err
	}

	// Create the token encoder.
	th, err := NewTokenHandler(sc.EncKey)
	if err != nil {
		return err
	}

	// Get the tls configuration. If the rootCA.crt file exists, that will
	// be used. Otherwise, use the system's certs.
	config, err := loadTlsConfig()
	if err != nil {
		return err
	}

	// Listen on the ListenAddr.
	ln, err := tls.Listen("tcp", sc.ListenAddr, config)
	if err != nil {
		return err
	}

	// Accept connections forever.
	for {
		rConn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleServerConnection(rConn, th)
	}

	return nil
}
