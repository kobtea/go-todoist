package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var (
	projectColor  string
	projectIndent string
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
		project := todoist.Project{
			Name: name,
		}
		if len(projectColor) > 0 {
			color, err := strconv.Atoi(projectColor)
			if err != nil {
				return fmt.Errorf("Invalid project color: %s", projectColor)
			}
			project.Color = color
		}
		if len(projectIndent) > 0 {
			i, err := strconv.Atoi(projectIndent)
			if err != nil {
				return fmt.Errorf("Invalid project indent: %s", projectIndent)
			}
			project.Indent = i
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

var projectUpdateCmd = &cobra.Command{
	Use:   "update id [new_name]",
	Short: "update projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Require project ID to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("Invalid ID: %s", args[0])
		}
		client, err := newClient()
		if err != nil {
			return err
		}
		project := client.Project.Resolve(id)
		if project == nil {
			return fmt.Errorf("No such project id: %s", id)
		}
		if len(args) > 1 {
			project.Name = strings.Join(args[1:], " ")
		}
		if len(projectColor) > 0 {
			color, err := strconv.Atoi(projectColor)
			if err != nil {
				return fmt.Errorf("Invalid project color: %s", projectColor)
			}
			project.Color = color
		}
		if len(projectIndent) > 0 {
			i, err := strconv.Atoi(projectIndent)
			if err != nil {
				return fmt.Errorf("Invalid project indent: %s", projectIndent)
			}
			project.Indent = i
		}
		if _, err = client.Project.Update(*project); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		syncedProject := client.Project.Resolve(id)
		if syncedProject == nil {
			return errors.New("Failed to add this project. It may be failed to sync.")
		}
		fmt.Println("Successful updating project.")
		fmt.Println(util.ProjectTableString([]todoist.Project{*syncedProject}))
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete id ...",
	Short: "delete projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Require project ID to delete")
		}
		var ids []todoist.ID
		for _, i := range args {
			id, err := todoist.NewID(i)
			if err != nil {
				return err
			}
			ids = append(ids, id)
		}
		client, err := newClient()
		if err != nil {
			return err
		}
		if err = client.Project.Delete(ids); err != nil {
			return err
		}
		ctx := context.Background()
		if err = client.Commit(ctx); err != nil {
			return err
		}
		if err = client.FullSync(ctx, []todoist.Command{}); err != nil {
			return err
		}
		fmt.Println("Successful deleting of project(s).")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectAddCmd.Flags().StringVarP(&projectColor, "color", "c", "7", "color")
	projectAddCmd.Flags().StringVarP(&projectIndent, "indent", "i", "1", "indent")
	projectCmd.AddCommand(projectAddCmd)
	projectUpdateCmd.Flags().StringVarP(&projectColor, "color", "c", "", "color")
	projectUpdateCmd.Flags().StringVarP(&projectIndent, "indent", "i", "", "indent")
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
}
