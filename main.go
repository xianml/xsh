package main

import (
	"fmt"
	"os"

	"github.com/xian/xsh/internal/config"
	"github.com/xian/xsh/internal/shell"
)

func main() {
	if os.Getenv("XSH_SESSION") == "true" {
		fmt.Fprintln(os.Stderr, "Error: Nested xsh sessions are not supported.")
		os.Exit(1)
	}
	os.Setenv("XSH_SESSION", "true")

	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("xsh v0.0.1 - AI Powered Shell")
		return
	}

	cfg := config.Load()
	xshell, err := shell.NewShell(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating shell: %v\n", err)
		os.Exit(1)
	}
	defer xshell.Goodbye()

	if err := xshell.Run(); err != nil {
		if err.Error() != "EOF" {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}
