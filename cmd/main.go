package main

import (
	"os"

	"github.com/faja/tornimo-agent/cmd/agent"
)

func main() {
	if err := agent.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
