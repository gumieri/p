package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of p",
	Run: func(cmd *cobra.Command, args []string) {
		typist.Println(Version)
	},
}
