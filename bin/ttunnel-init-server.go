package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"os"
	"strconv"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf("Usage: %v <host_name> <key_bits>\n", os.Args[0])
	fmt.Println("")
	fmt.Println("host_name:")
	fmt.Println("    The server's public host name.")
	fmt.Println("")
	fmt.Println("key_bits:")
	fmt.Println("    The number of bits to use for the RSA key pair.")
	fmt.Println("    Valid values are 2048, 3072, 4096, etc.")
	fmt.Println("")
}

func main() {
	if len(os.Args) != 3 {
		printUsage()
		return
	}

	// Read in command line arguments.
	host := os.Args[1]
	rsaBits, err := strconv.ParseInt(os.Args[2], 10, 32)
	if err != nil {
		printUsage()
		return
	}

	if err = tt.InitServer(host, int(rsaBits)); err != nil {
		printUsage()
		return
	}
}
