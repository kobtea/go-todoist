package cmd

import (
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"sort"
)

// todayCmd represents the today command
var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "show today's tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		items := client.Item.FindByDueDate(todoist.Today())
		sort.Slice(items, func(i, j int) bool {
			return items[i].DueDateUtc.Before(items[j].DueDateUtc)
		})
		relations := client.Relation.Items(items)
		fmt.Println(util.ItemTableString(items, relations, func(i todoist.Item) todoist.Time { return i.DueDateUtc }))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(todayCmd)
}
