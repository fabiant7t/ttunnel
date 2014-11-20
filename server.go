package ttunnel

import (
	"crypto/tls"
	"log"
	"net"
)

func RunServer(listenAddr string) (err error) {

	// Load the key and certificate.
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return
	}

	// Create TLS configuration for the server.
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	// Accept connections on the given address.
	ln, err := tls.Listen("tcp", listenAddr, &config)
	if err != nil {
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go serverHandler(conn)
	}

	return
}

func serverHandler(conn net.Conn) {
	cc := ClientConfig{}
	var pwd []byte
	var cConn net.Conn

	// Get the client's name.
	name, err := readBytes(conn, 256)
	if err != nil {
		log.Printf("Error reading name:\n    %v\n", err)
		goto closeConn
	}

	// Load the client's configuration.
	if err = cc.Load(clientsPath(string(name))); err != nil {
		log.Printf("Failed to load client's configuration:\n    %v\n", err)
		goto closeConn
	}

	// Get the client's password.
	pwd, err = readBytes(conn, 256)
	if err != nil {
		log.Printf("Error reading password:\n    %v\n", err)
		goto closeConn
	}

	// Check the client's password.
	if !cc.PwdMatches(pwd) {
		log.Printf("Invalid password for client %v:\n    %v\n", name, err)
		goto closeConn
	}

	// Connect to the remote address for the client.
	cConn, err = net.Dial("tcp", cc.ConnectAddr)
	if err != nil {
		log.Printf("Failed to connect to address %v for client %v:\n    %v\n",
			cc.ConnectAddr, name, err)
		goto closeConn
	}

	// Forward traffic.
	log.Printf("Forwarding traffic for client %v.\n", string(name))
	copyDuplex(cConn, conn)
	return

closeConn:
	conn.Close()
}
