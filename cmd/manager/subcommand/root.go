package subcommand

import (
	"fmt"

	"github.com/durandj/ley/internal/manager"
	"github.com/durandj/ley/internal/manager/configuration"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a root command for the CLI to run.
func NewRootCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "manager",
		Short: "Ley network manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			config, err := configuration.NewFromEnvironment()
			if err != nil {
				return fmt.Errorf("Unable to load service configuration: %w", err)
			}

			server, err := manager.New(config)
			if err != nil {
				return fmt.Errorf("Unable to setup server: %w", err)
			}

			return server.Run(ctx)
		},
	}

	return &cmd
}
