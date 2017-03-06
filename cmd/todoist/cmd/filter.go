package cmd

import (
	"fmt"

	"context"
	"errors"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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

var filterAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return err
		}
		filter := todoist.Filter{
			Name:  name,
			Query: query,
		}
		colorStr, err := cmd.Flags().GetString("color")
		if err != nil {
			return errors.New("Invalid filter color")
		}
		if len(colorStr) > 0 {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("Invalid filter color: %s", colorStr)
			}
			filter.Color = color
		}
		if _, err = client.Filter.Add(filter); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		filters := client.Filter.FindByName(name)
		if len(filters) == 0 {
			return errors.New("Failed to add this filter. It may be failed to sync.")
		}
		syncedFilter := filters[len(filters)-1]
		fmt.Println("Successful addition of a filter.")
		fmt.Println(util.FilterTableString([]todoist.Filter{syncedFilter}))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(filterCmd)
	filterCmd.AddCommand(filterListCmd)
	filterAddCmd.Flags().StringP("query", "q", "", "query")
	filterAddCmd.Flags().StringP("color", "c", "12", "color")
	filterCmd.AddCommand(filterAddCmd)
}
