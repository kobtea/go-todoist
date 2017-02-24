package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"strings"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "subcommand for project",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "list projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		projects := client.Project.GetAll()
		fmt.Println(util.ProjectTableString(projects))
		return nil
	},
}

var projectAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		color, err := cmd.Flags().GetInt("color")
		if err != nil {
			return err
		}
		indent, err := cmd.Flags().GetInt("indent")
		if err != nil {
			return err
		}
		project := todoist.Project{
			Name:   name,
			Color:  color,
			Indent: indent,
		}
		if _, err = client.Project.Add(project); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		projects := client.Project.FindByName(name)
		if len(projects) == 0 {
			return errors.New("Failed to add this project. It may be failed to sync.")
		}
		// it may not be new project
		syncedProject := projects[len(projects)-1]
		fmt.Println("Successful addition of an project.")
		fmt.Println(util.ProjectTableString([]todoist.Project{syncedProject}))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectAddCmd.Flags().IntP("color", "c", 7, "color")
	projectAddCmd.Flags().IntP("indent", "i", 1, "indent")
	projectCmd.AddCommand(projectListCmd, projectAddCmd)
}
