package server

import (
	"github.com/spf13/cobra"
)

type app struct {
	*cobra.Command
}

//nolint:revive
func NewApp() *app {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "auth server",
	}
	app := &app{Command: cmd}
	app.RunE = func(c *cobra.Command, args []string) error {
		return app.run()
	}
	return app
}
