package main

import (
	"flag"
	"fmt"
)

var (
	dir     = flag.String("d", ".", "Root directory of your Go application")
	verbose = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	fmt.Printf("Your dir: %s\n", *dir)
}
