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

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "subcommand for project",
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "list projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
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
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		project := todoist.Project{
			Name: name,
		}
		colorStr, err := cmd.Flags().GetString("color")
		if err != nil {
			return errors.New("Invalid project color")
		}
		if len(colorStr) > 0 {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("Invalid project color: %s", colorStr)
			}
			project.Color = color
		}
		indentStr, err := cmd.Flags().GetString("indent")
		if err != nil {
			return errors.New("Invalid project indent")
		}
		if len(indentStr) > 0 {
			i, err := strconv.Atoi(indentStr)
			if err != nil {
				return fmt.Errorf("Invalid project indent: %s", indentStr)
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
		client, err := util.NewClient()
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
		colorStr, err := cmd.Flags().GetString("color")
		if err != nil {
			return errors.New("Invalid project color")
		}
		if len(colorStr) > 0 {
			color, err := strconv.Atoi(colorStr)
			if err != nil {
				return fmt.Errorf("Invalid project color: %s", colorStr)
			}
			project.Color = color
		}
		indentStr, err := cmd.Flags().GetString("indent")
		if err != nil {
			return errors.New("Invalid project indent")
		}
		if len(indentStr) > 0 {
			i, err := strconv.Atoi(indentStr)
			if err != nil {
				return fmt.Errorf("Invalid project indent: %s", indentStr)
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
	Use:   "delete id [...]",
	Short: "delete projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(args, client.Project.Delete)
		}); err != nil {
			return err
		}
		fmt.Println("Successful deleting of project(s).")
		return nil
	},
}

var projectArchiveCmd = &cobra.Command{
	Use:   "archive id [...]",
	Short: "archive projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(args, client.Project.Archive)
		}); err != nil {
			return err
		}
		fmt.Println("Successful archiving of project(s).")
		return nil
	},
}

var projectUnarchiveCmd = &cobra.Command{
	Use:   "unarchive id [...]",
	Short: "unarchive projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			return util.ProcessIDs(args, client.Project.Unarchive)
		}); err != nil {
			return err
		}
		fmt.Println("Successful un-archiving of project(s).")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectAddCmd.Flags().StringP("color", "c", "7", "color")
	projectAddCmd.Flags().StringP("indent", "i", "1", "indent")
	projectCmd.AddCommand(projectAddCmd)
	projectUpdateCmd.Flags().StringP("color", "c", "", "color")
	projectUpdateCmd.Flags().StringP("indent", "i", "", "indent")
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectUnarchiveCmd)
}
