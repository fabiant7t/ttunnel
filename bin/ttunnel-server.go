package main

import (
	tt "github.com/johnnylee/ttunnel"
	"log"
)

func main() {
	if err := tt.RunServer(); err != nil {
		log.Printf("Error running server:\n    %v\n", err)
	}
}
