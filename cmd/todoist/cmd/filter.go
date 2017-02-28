package cmd

import (
	"fmt"

	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/spf13/cobra"
)

// filterCmd represents the filter command
var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "subcommand for filter",
}

var filterListCmd = &cobra.Command{
	Use:   "list",
	Short: "list filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		filters := client.Filter.GetAll()
		fmt.Println(util.FilterTableString(filters))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(filterCmd)
	filterCmd.AddCommand(filterListCmd)
}
