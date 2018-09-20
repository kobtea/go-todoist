package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generate completion script",
}

var completionBashCmd = &cobra.Command{
	Use: "bash",
	Short: "generate bash completion script",
	Run: func(cmd *cobra.Command, args []string) {
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.AddCommand(completionBashCmd)
}
