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

var filterUpdateCmd = &cobra.Command{
	Use:   "update id [new_name]",
	Short: "update filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Require filter ID to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("Invalid ID: %s", args[0])
		}
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		filter := client.Filter.Resolve(id)
		if filter == nil {
			return fmt.Errorf("No such filter id: %s", id)
		}
		if len(args) > 1 {
			filter.Name = strings.Join(args[1:], " ")
		}
		query, err := cmd.Flags().GetString("query")
		if err != nil {
			return err
		}
		if len(query) > 0 {
			filter.Query = query
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
		if _, err = client.Filter.Update(*filter); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		syncedFilter := client.Filter.Resolve(id)
		if syncedFilter == nil {
			return errors.New("Failed to add this filter. It may be failed to sync.")
		}
		fmt.Println("Successful updating filter.")
		fmt.Println(util.FilterTableString([]todoist.Filter{*syncedFilter}))
		return nil
	},
}

var filterDeleteCmd = &cobra.Command{
	Use:   "delete id [...]",
	Short: "delete filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(
				args,
				func(ids []todoist.ID) error {
					for _, id := range ids {
						if err := client.Filter.Delete(id); err != nil {
							return err
						}
					}
					return nil
				})
		}); err != nil {
			return err
		}
		fmt.Println("Successful deleting of filter(s).")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(filterCmd)
	filterCmd.AddCommand(filterListCmd)
	filterAddCmd.Flags().StringP("query", "q", "", "query")
	filterAddCmd.Flags().StringP("color", "c", "12", "color")
	filterCmd.AddCommand(filterAddCmd)
	filterUpdateCmd.Flags().StringP("query", "q", "", "query")
	filterUpdateCmd.Flags().StringP("color", "c", "", "color")
	filterCmd.AddCommand(filterUpdateCmd)
	filterCmd.AddCommand(filterDeleteCmd)
}
