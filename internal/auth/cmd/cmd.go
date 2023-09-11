package cmd

import (
	"github.com/and-period/furumane/internal/auth/cmd/server"
	"github.com/spf13/cobra"
)

func RegisterCommand(registry *cobra.Command) {
	registry.AddCommand(server.NewApp().Command)
}
