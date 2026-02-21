package main

import (
	"fmt"
	"os"

	"github.com/jgfranco17/llm-workspace/internal/cli"
)

func main() {
	r := cli.New()
	if err := r.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to run tools: %s\n", err.Error())
		os.Exit(1)
	}
}
