/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/drycc/workflow-cli/cmd"
)

func main() {
	rootCmd := cmd.NewDryccCommand()

	// Get config file path
	config := "~/.drycc/client.json"
	if v, ok := os.LookupEnv("DRYCC_PROFILE"); ok {
		config = v
	}

	if err := cmd.ExecuteWithPlugins(rootCmd, config); err != nil {
		os.Exit(1)
	}
}
