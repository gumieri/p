package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of p",
	Run: func(cmd *cobra.Command, args []string) {
		debug.SetGCPercent(-1)
		t.Outln(Version)
	},
}
