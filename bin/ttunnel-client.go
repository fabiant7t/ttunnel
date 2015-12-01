package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
)

func printUsage() {
	fmt.Println("")
	fmt.Println("Usage: %v")
	fmt.Println("")
}

func main() {
	if err := tt.RunClient(); err != nil {
		fmt.Printf("Error running client:\n    %v\n", err)
	}
}
