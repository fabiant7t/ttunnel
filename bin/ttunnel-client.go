package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"os"
	"runtime"
)

func printUsage() {
	fmt.Println("")
	fmt.Printf("Usage: %v <config_name>\n", os.Args[0])
	fmt.Println("")
	fmt.Println("config_name:")
	fmt.Println("    The name of the configuration to use.")
	fmt.Println("")
}

func main() {
	if len(os.Args) != 2 {
		printUsage()
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := tt.RunClient(os.Args[1]); err != nil {
		fmt.Printf("Error running client:\n    %v\n", err)
	}
}
