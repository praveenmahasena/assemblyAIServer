package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/praveenmahasena/server/internal"
)

func main() {
	if err := internal.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		syscall.Exit(255)
	}
}
