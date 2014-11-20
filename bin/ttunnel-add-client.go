package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"os"
	"strconv"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf(
		"Usage: %v <name> <host_addr> <connect_addr> <local_port>\n",
		os.Args[0])
	fmt.Println("")
	fmt.Println("name:")
	fmt.Println("    A unique name for the user's tunnel. This must be")
	fmt.Println("    a valid filename, and must be unique on the server")
	fmt.Println("    and client.")
	fmt.Println("")
	fmt.Println("host_addr:")
	fmt.Println("    The address of the tunnel server.")
	fmt.Println("    This has the format <address>:<port>.")
	fmt.Println("")
	fmt.Println("connect_addr:")
	fmt.Println("    The address to forward the client's traffic to.")
	fmt.Println("    This has the format <address>:<port>.")
	fmt.Println("")
	fmt.Println("local_port:")
	fmt.Println("    The port on the client machine on which to accept")
	fmt.Println("    connections.")
	fmt.Println("")
}

func main() {
	if len(os.Args) != 5 {
		printUsage()
		return
	}

	name := os.Args[1]
	hostAddr := os.Args[2]
	connectAddr := os.Args[3]
	localPort, err := strconv.ParseInt(os.Args[4], 10, 32)
	if err != nil {
		printUsage()
		return
	}

	err = tt.AddClient(name, hostAddr, connectAddr, int32(localPort))
	if err != nil {
		printUsage()
		return
	}
}
