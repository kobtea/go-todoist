package cmd

import (
	"fmt"

	"errors"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
)

// inboxCmd represents the inbox command
var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "show inbox tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		projects := client.Project.FindByName("Inbox")
		if len(projects) != 1 {
			errors.New("Failed to detect inbox. It may exist multiple inbox.")
		}
		inbox := projects[0]
		items := client.Item.FindByProjectIDs([]todoist.ID{inbox.ID})
		relations := client.Relation.Items(items)
		fmt.Println(util.ItemTableString(items, relations))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(inboxCmd)
}
