package cmd

import (
	"fmt"

	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/spf13/cobra"
)

// itemCmd represents the item command
var itemCmd = &cobra.Command{
	Use:   "item",
	Short: "subcommand for item",
}

var itemListCmd = &cobra.Command{
	Use:   "list",
	Short: "list items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		items := client.Item.GetAll()
		relations := client.Relation.Items(items)
		fmt.Println(util.ItemTableString(items, relations))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(itemCmd)
	itemCmd.AddCommand(itemListCmd)
}
