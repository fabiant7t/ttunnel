package ttunnel

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// forwardConnections fowards traffic between the two connections.
func forwardConnections(c1, c2 net.Conn) {
	go copyConn(c1, c2)
	go copyConn(c2, c1)
}

func copyConn(in, out net.Conn) {
	var n int
	var err error

	buf := make([]byte, 1048576)

	for {
		n, err = in.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading:\n    %v\n", err)
			}
			in.Close()
			out.Close()
			return
		}

		_, err = out.Write(buf[:n])
		if err != nil {
			if err != io.EOF {
				log.Printf("Error writing:\n    %v\n", err)
			}
			in.Close()
			out.Close()
			return
		}
	}
}

func clientHandler(lConn net.Conn, cc *ClientConfig) {
	// Connect to server.
	rConn, err := tls.Dial("tcp", cc.Host, nil)
	if err != nil {
		log.Printf("Error dialing server:\n    %v\n", err)
		lConn.Close()
		return
	}

	// Send the token.
	err = binary.Write(rConn, binary.LittleEndian, uint16(len(cc.Token)))
	if err != nil {
		log.Printf("Error sending token length:\n    %v\n", err)
		rConn.Close()
		lConn.Close()
		return
	}

	_, err = rConn.Write(cc.Token)
	if err != nil {
		log.Printf("Error sending token:\n    %v\n", err)
		rConn.Close()
		lConn.Close()
		return
	}

	// Forward connections.
	forwardConnections(lConn, rConn)
}

// RunClient
func RunClient(tunnelName string) error {

	// Load the client's configuration.
	cc, err := ReadClientConfig(tunnelName)
	if err != nil {
		return err
	}

	// Accept connections on the local port.
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cc.Port))
	if err != nil {
		return err
	}

	for {
		lConn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection:\n     %v\n", err)
			continue
		}

		go clientHandler(lConn, &cc)
	}

	return nil
}
