package cmd

import (
	"fmt"

	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/spf13/cobra"
)

// labelCmd represents the label command
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "subcommand for label",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "list projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		labels := client.Label.GetAll()
		fmt.Println(util.LabelTableString(labels))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
}
