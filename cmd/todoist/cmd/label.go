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
	"strconv"
	"strings"
)

// labelCmd represents the label command
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "subcommand for label",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "list label",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		labels := client.Label.GetAll()
		fmt.Println(util.LabelTableString(labels))
		return nil
	},
}

var labelAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add label",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		label := todoist.Label{
			Name: name,
		}
		colorStr, err := cmd.Flags().GetString("color")
		if err != nil {
			return errors.New("Invalid label color")
		}
		if len(colorStr) > 0 {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("Invalid label color: %s", colorStr)
			}
			label.Color = color
		}
		if _, err = client.Label.Add(label); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		labels := client.Label.FindByName(name)
		if len(labels) == 0 {
			return errors.New("Failed to add this label. It may be failed to sync.")
		}
		// it may not be new label
		syncedLabel := labels[len(labels)-1]
		fmt.Println("Successful addition of a label.")
		fmt.Println(util.LabelTableString([]todoist.Label{syncedLabel}))
		return nil
	},
}

var labelUpdateCmd = &cobra.Command{
	Use:   "update id [new_name]",
	Short: "update label",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Require label ID to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("Invalid ID: %s", args[0])
		}
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		label := client.Label.Resolve(id)
		if label == nil {
			return fmt.Errorf("No such label id: %s", id)
		}
		if len(args) > 1 {
			label.Name = strings.Join(args[1:], " ")
		}
		colorStr, err := cmd.Flags().GetString("color")
		if err != nil {
			return errors.New("Invalid label color")
		}
		if len(colorStr) > 0 {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("Invalid label color: %s", colorStr)
			}
			label.Color = color
		}
		if _, err = client.Label.Update(*label); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		syncedLabel := client.Label.Resolve(id)
		if syncedLabel == nil {
			return errors.New("Failed to add this label. It may be failed to sync.")
		}
		fmt.Println("Successful updating label.")
		fmt.Println(util.LabelTableString([]todoist.Label{*syncedLabel}))
		return nil
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete id [...]",
	Short: "delete labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(
				args,
				func(ids []todoist.ID) error {
					var labels []todoist.Label
					for _, id := range ids {
						label := client.Label.Resolve(id)
						if label == nil {
							return fmt.Errorf("invalid id: %s", id)
						}
						labels = append(labels, *label)
					}
					fmt.Println(util.LabelTableString(labels))

					reader := bufio.NewReader(os.Stdin)
					fmt.Print("are you sure to delete above label(s)? (y/[n]): ")
					ans, err := reader.ReadString('\n')
					if ans != "y\n" || err != nil {
						fmt.Println("abort")
						return errors.New("abort")
					}
					for _, id := range ids {
						if err := client.Label.Delete(id); err != nil {
							return err
						}
					}
					return nil
				})
		}); err != nil {
			if err.Error() == "abort" {
				return nil
			}
			return err
		}
		fmt.Println("Successful deleting of label(s).")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
	labelAddCmd.Flags().StringP("color", "c", "7", "color")
	labelCmd.AddCommand(labelAddCmd)
	labelUpdateCmd.Flags().StringP("color", "c", "", "color")
	labelCmd.AddCommand(labelUpdateCmd)
	labelCmd.AddCommand(labelDeleteCmd)
}
