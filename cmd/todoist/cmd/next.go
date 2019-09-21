package cmd

import (
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"sort"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "show next 7 days tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		var items []todoist.Item
		for _, i := range client.Item.FindByDueDate(todoist.Next7Days()) {
			if !i.IsChecked() {
				items = append(items, i)
			}
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Due.Date.Before(items[j].Due.Date)
		})
		relations := client.Relation.Items(items)
		fmt.Println(util.ItemTableString(items, relations, func(i todoist.Item) todoist.Time { return i.Due.Date }))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(nextCmd)
}
