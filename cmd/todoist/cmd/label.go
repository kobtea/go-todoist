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
	Use:   "add [name]",
	Short: "add label",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		if len(name) == 0 {
			return errors.New("require label name")
		}
		opts := todoist.NewLabelOpts{}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			opts.Color = color
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			opts.ItemOrder = order
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			opts.IsFavorite = todoist.IntBool(favorite)
		}
		label, err := todoist.NewLabel(name, &opts)
		if label == nil {
			return errors.New("failed to initialize a label")
		}
		if err != nil {
			return err
		}
		if _, err := client.Label.Add(*label); err != nil {
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
			return errors.New("failed to add this label. it may be failed to sync")
		}
		// it may not be new label
		syncedLabel := labels[len(labels)-1]
		fmt.Println("succeeded to add a label")
		fmt.Println(util.LabelTableString([]todoist.Label{syncedLabel}))
		return nil
	},
}

var labelUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "update label",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("require label id to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		label := client.Label.Resolve(id)
		if label == nil {
			return fmt.Errorf("no such label id: %s", id)
		}
		if name, err := cmd.Flags().GetString("name"); err != nil {
			return err
		} else {
			if len(name) != 0 {
				label.Name = name
			}
		}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("color") {
				label.Color = color
			}
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("order") {
				label.ItemOrder = order
			}
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			return err
		} else {
			if favorite {
				label.IsFavorite = true
			}
		}
		if unFavorite, err := cmd.Flags().GetBool("un-favorite"); err != nil {
			return err
		} else {
			if unFavorite {
				label.IsFavorite = false
			}
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
			return errors.New("failed to add this label. it may be failed to sync")
		}
		fmt.Println("succeeded to update the label")
		fmt.Println(util.LabelTableString([]todoist.Label{*syncedLabel}))
		return nil
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete label",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			if len(args) == 0 {
				return errors.New("require label id to delete")
			}
			return util.ProcessID(args[0], func(id todoist.ID) error {
				label := client.Label.Resolve(id)
				if label == nil {
					return fmt.Errorf("invalid label id: %s", id)
				}
				fmt.Println(util.LabelTableString([]todoist.Label{*label}))
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("are you sure to delete above label? (y/[n]): ")
				ans, err := reader.ReadString('\n')
				if ans != "y\n" || err != nil {
					fmt.Println("abort")
					return errors.New("abort")
				}
				return client.Label.Delete(id)
			})
		}); err != nil {
			if err.Error() == "abort" {
				return nil
			}
			return err
		}
		fmt.Println("succeeded to delete the label")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
	labelAddCmd.Flags().IntP("color", "c", 47, "color")
	labelAddCmd.Flags().Int("order", 0, "item order")
	labelAddCmd.Flags().Bool("favorite", false, "is favorite")
	labelCmd.AddCommand(labelAddCmd)
	labelUpdateCmd.Flags().String("name", "", "name of the label")
	labelUpdateCmd.Flags().IntP("color", "c", 47, "color")
	labelUpdateCmd.Flags().Int("order", 0, "item order")
	labelUpdateCmd.Flags().Bool("favorite", false, "is favorite")
	labelUpdateCmd.Flags().Bool("un-favorite", false, "is not favorite")
	labelCmd.AddCommand(labelUpdateCmd)
	labelCmd.AddCommand(labelDeleteCmd)
}
