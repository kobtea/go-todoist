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
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		items := client.Item.GetAll()
		relations := client.Relation.Items(items)
		fmt.Println(util.ItemTableString(items, relations, func(i todoist.Item) todoist.Time { return i.DueDateUtc }))
		return nil
	},
}

var itemAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		content := strings.Join(args, " ")
		item := todoist.Item{Content: content}

		projectName, err := cmd.Flags().GetString("project")
		if err != nil {
			return errors.New("invalid project name")
		}
		if project := client.Project.FindOneByName(projectName); project != nil {
			item.ProjectID = project.ID
		}

		labelNames, err := cmd.Flags().GetString("label")
		if err != nil {
			return errors.New("invalid label name")
		}
		for _, labelName := range strings.Split(labelNames, ",") {
			if label := client.Label.FindOneByName(labelName); label != nil {
				item.Labels = append(item.Labels, label.ID)
			}
		}

		priority, err := cmd.Flags().GetInt("priority")
		if err != nil {
			return errors.New("invalid priority")
		}
		item.Priority = priority

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
		fmt.Println("Successful addition of an item.")
		fmt.Println(util.ItemTableString([]todoist.Item{syncedItem}, relations, func(i todoist.Item) todoist.Time { return i.DueDateUtc }))
		return nil
	},
}

var itemDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete items",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(
				args,
				func(ids []todoist.ID) error {
					return client.Item.Delete(ids)
				})
		}); err != nil {
			return err
		}
		fmt.Println("Successful deleting of item(s).")
		return nil
	},
}

var itemMoveCmd = &cobra.Command{
	Use:   "move",
	Short: "move the project of the item",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		if len(args) < 1 {
			return errors.New("Require item ID to move")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("Invalid ID: %s", args[0])
		}
		item := client.Item.Resolve(id)
		if item == nil {
			return fmt.Errorf("No such item id: %s", id)
		}
		pidstr, err := cmd.Flags().GetString("project")
		if err != nil {
			return errors.New("Invalid project id")
		}
		if len(pidstr) == 0 {
			return errors.New("Require project ID to move")
		}
		pid, err := todoist.NewID(pidstr)
		if err != nil {
			return fmt.Errorf("Invalid ID: %s", args[0])
		}
		if err = client.Item.Move(map[todoist.ID][]todoist.ID{item.ProjectID: {item.ID}}, pid); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		syncedItem := client.Item.Resolve(id)
		if syncedItem == nil {
			return errors.New("Failed to move this item. It may be failed to sync.")
		}
		relations := client.Relation.Items([]todoist.Item{*syncedItem})
		fmt.Println("Successful move item.")
		fmt.Println(util.ItemTableString([]todoist.Item{*syncedItem}, relations, func(i todoist.Item) todoist.Time { return i.DueDateUtc }))
		return nil
	},
}

var itemCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "complete items",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(
				args,
				func(ids []todoist.ID) error {
					return client.Item.Complete(ids, true)
				})
		}); err != nil {
			return err
		}
		fmt.Println("Successful completion of item(s).")
		return nil
	},
}

var itemUncompleteCmd = &cobra.Command{
	Use:   "uncomplete",
	Short: "uncomplete items",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(
				args,
				func(ids []todoist.ID) error {
					restoreState := map[todoist.ID][]string{}
					return client.Item.Uncomplete(ids, true, restoreState)
				})
		}); err != nil {
			return err
		}
		fmt.Println("Successful uncompletion of item(s).")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(itemCmd)
	itemCmd.AddCommand(itemListCmd)
	itemAddCmd.Flags().StringP("project", "p", "inbox", "project name")
	itemAddCmd.Flags().StringP("label", "l", "", "label name(s) (delimiter: ,)")
	// itemAddCmd.Flags().StringP("due", "d", "", "due date")
	itemAddCmd.Flags().Int("priority", 1, "priority")
	itemCmd.AddCommand(itemAddCmd)
	itemCmd.AddCommand(itemDeleteCmd)
	itemMoveCmd.Flags().StringP("project", "p", "", "project")
	itemCmd.AddCommand(itemMoveCmd)
	itemCmd.AddCommand(itemCompleteCmd)
	itemCmd.AddCommand(itemUncompleteCmd)
}
