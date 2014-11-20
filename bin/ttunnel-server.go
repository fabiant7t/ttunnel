package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"os"
	"runtime"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf("Usage: %v <listen_addr>\n", os.Args[0])
	fmt.Println("")
	fmt.Println("listen_addr:")
	fmt.Println("    The address to listen on. This must be of the form")
	fmt.Println("    <address.:<port>.")
	fmt.Println("")
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := tt.RunServer(os.Args[1]); err != nil {
		fmt.Printf("Error running server:\n    %v\n", err)
	}
}
