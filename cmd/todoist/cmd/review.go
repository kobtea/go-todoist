package cmd

import (
	"errors"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"sort"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "show completed items",
	RunE: func(cmd *cobra.Command, args []string) error {
		sortKey, err := cmd.Flags().GetString("sort")
		if err != nil {
			return errors.New("Invalid sort key")
		}
		if sortKey != "time" && sortKey != "project" {
			fmt.Println("Invalid sort key was given. Use default key (time).")
		}

		client, err := util.NewClient()
		if err != nil {
			return err
		}
		completed, err := client.Completed.GetAll()
		if err != nil {
			return err
		}
		// first sort with time
		sort.Slice(completed.Items, func(i, j int) bool {
			return completed.Items[i].CompletedDate.Before(completed.Items[j].CompletedDate)
		})
		sort.Slice(completed.Items, func(i, j int) bool {
			switch sortKey {
			case "project":
				return completed.Items[i].ProjectID < completed.Items[j].ProjectID
			default:
				return false
			}
		})
		relations := client.Relation.Items(completed.Items)
		fmt.Println(util.ItemTableString(completed.Items, relations, func(i todoist.Item) todoist.Time { return i.CompletedDate }))
		return nil
	},
}

func init() {
	reviewCmd.Flags().StringP("sort", "s", "time", "sort key")
	RootCmd.AddCommand(reviewCmd)
}
