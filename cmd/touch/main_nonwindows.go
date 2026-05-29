//go:build !windows

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "touch.exe is a Windows-first application. Build with GOOS=windows.")
	os.Exit(1)
}
