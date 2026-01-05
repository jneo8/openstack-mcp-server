package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print the version, commit hash, and build date of the OpenStack MCP server.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("OpenStack MCP Server\n")
			fmt.Printf("  Version:    %s\n", version)
			fmt.Printf("  Commit:     %s\n", commit)
			fmt.Printf("  Build Date: %s\n", buildDate)
		},
	}
}
