package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lateralusd/tuid/fetcher"
)

func main() {
	userPath := flag.String("users", "users", "read users from file")
	flag.Parse()
	f := fetcher.NewFetcher()
	if err := f.Monitor(*userPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error ocurred: %v", err)
		os.Exit(1)
	}
}
