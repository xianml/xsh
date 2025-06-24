package main

import (
	"fmt"
	"os"

	"github.com/xian/xsh/internal/config"
	"github.com/xian/xsh/internal/shell"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("xsh v1.0.0 - AI Powered Shell")
		return
	}

	cfg := config.Load()
	xshell, err := shell.NewShell(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating shell: %v\n", err)
		os.Exit(1)
	}

	if err := xshell.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
