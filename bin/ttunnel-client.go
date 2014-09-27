package main

import (
	"fmt"
	tt "github.com/johnnylee/ttunnel"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <tunnel-name>\n", os.Args[0])
		return
	}

	if err := tt.RunClient(os.Args[1]); err != nil {
		log.Printf("Error running client:\n    %v\n", err)
	}
}
