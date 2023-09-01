package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("testing linter")
	os.Exit(0) // want "os.Exit in main module"
}
