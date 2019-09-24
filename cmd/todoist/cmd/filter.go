package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"os"
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
	Use:   "add [name]",
	Short: "add filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		if len(name) == 0 {
			return errors.New("require filter name")
		}
		var query string
		if query, err = cmd.Flags().GetString("query"); err != nil {
			return err
		} else {
			if len(query) == 0 {
				return errors.New("require filter query")
			}
		}
		opts := todoist.NewFilterOpts{}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			opts.Color = color
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			return err
		} else {
			opts.ItemOrder = order
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			return err
		} else {
			opts.IsFavorite = todoist.IntBool(favorite)
		}
		filter, err := todoist.NewFilter(name, query, &opts)
		if filter == nil {
			return errors.New("failed to initialize a filter")
		}
		if err != nil {
			return err
		}
		if _, err = client.Filter.Add(*filter); err != nil {
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
		fmt.Println("succeeded to add a filter")
		fmt.Println(util.FilterTableString([]todoist.Filter{syncedFilter}))
		return nil
	},
}

var filterUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "update filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("require filter id to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		filter := client.Filter.Resolve(id)
		if filter == nil {
			return fmt.Errorf("no such filter id: %s", id)
		}
		if name, err := cmd.Flags().GetString("name"); err != nil {
			return err
		} else {
			if len(name) != 0 {
				filter.Name = name
			}
		}
		if query, err := cmd.Flags().GetString("query"); err != nil {
			return err
		} else {
			if len(query) != 0 {
				filter.Query = query
			}
		}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("color") {
				filter.Color = color
			}
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("order") {
				filter.ItemOrder = order
			}
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			return err
		} else {
			if favorite {
				filter.IsFavorite = true
			}
		}
		if unFavorite, err := cmd.Flags().GetBool("un-favorite"); err != nil {
			return err
		} else {
			if unFavorite {
				filter.IsFavorite = false
			}
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
			return errors.New("failed to add this filter. it may be failed to sync")
		}
		fmt.Println("succeeded to update the filter")
		fmt.Println(util.FilterTableString([]todoist.Filter{*syncedFilter}))
		return nil
	},
}

var filterDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete filter",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			if len(args) == 0 {
				return errors.New("require filter id to delete")
			}
			return util.ProcessID(args[0], func(id todoist.ID) error {
				filter := client.Filter.Resolve(id)
				if filter == nil {
					return fmt.Errorf("invalid filter id: %s", id)
				}
				fmt.Println(util.FilterTableString([]todoist.Filter{*filter}))
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("are you sure to delete above filter? (y/[n]): ")
				ans, err := reader.ReadString('\n')
				if ans != "y\n" || err != nil {
					fmt.Println("abort")
					return errors.New("abort")
				}
				return client.Filter.Delete(id)
			})
		}); err != nil {
			if err.Error() == "abort" {
				return nil
			}
			return err
		}
		fmt.Println("succeeded to delete the filter")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(filterCmd)
	filterCmd.AddCommand(filterListCmd)
	filterAddCmd.Flags().StringP("query", "q", "", "query")
	filterAddCmd.Flags().IntP("color", "c", 47, "color")
	filterAddCmd.Flags().Int("order", 0, "item order")
	filterAddCmd.Flags().Bool("favorite", false, "is favorite")
	filterCmd.AddCommand(filterAddCmd)
	filterUpdateCmd.Flags().String("name", "", "name of the filter")
	filterUpdateCmd.Flags().StringP("query", "q", "", "query")
	filterUpdateCmd.Flags().IntP("color", "c", 47, "color")
	filterUpdateCmd.Flags().Int("order", 0, "item order")
	filterUpdateCmd.Flags().Bool("favorite", false, "is favorite")
	filterUpdateCmd.Flags().Bool("un-favorite", false, "is not favorite")
	filterCmd.AddCommand(filterUpdateCmd)
	filterCmd.AddCommand(filterDeleteCmd)
}
