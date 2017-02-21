package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"sort"
	"strings"
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

var itemAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		content := strings.Join(args, " ")
		item := todoist.Item{Content: content}
		if _, err = client.Item.Add(item); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		// retrieve the item
		items := client.Item.FindByContent(content)
		if len(items) == 0 {
			return errors.New("Failed to add this item. It may be failed to sync.")
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].DateAdded.Before(items[j].DateAdded)
		})
		syncedItem := items[len(items)-1]
		relations := client.Relation.Items([]todoist.Item{syncedItem})
		fmt.Println(util.ItemTableString([]todoist.Item{syncedItem}, relations))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(itemCmd)
	itemCmd.AddCommand(itemListCmd)
	itemCmd.AddCommand(itemAddCmd)
}
