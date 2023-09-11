package main

import (
	"os"

	"github.com/and-period/furumane/internal/auth/cmd"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{Use: "auth [command]"}
	cmd.RegisterCommand(c)
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
