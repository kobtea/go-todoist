package cmd

import (
	"context"
	"fmt"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncronize origin server",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		if err = client.FullSync(context.Background(), []todoist.Command{}); err != nil {
			return err
		}
		fmt.Printf("update sync token: %s", client.SyncToken)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
