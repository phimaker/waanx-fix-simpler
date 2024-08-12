package cmd

import (
	marketdata "github.com/phimaker/waanx-fix-simpler/cmd/market-data"
	"github.com/phimaker/waanx-fix-simpler/internal/version"
	"github.com/spf13/cobra"
)

var (
	// versionF flag prints the version and exits.
	versionF bool
)

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() error {

	c := &cobra.Command{
		Use: "waanx-adapter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if versionF {
				version.PrintVersion()
				return nil
			}
			return cmd.Usage()
		},
	}

	c.Flags().BoolVarP(&versionF, "version", "v", false, "show the version and exit")

	c.AddCommand(marketdata.Cmd)

	return c.Execute()
}
