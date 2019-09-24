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
	Use:   "add [name]",
	Short: "add a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		name := strings.Join(args, " ")
		if len(name) == 0 {
			return errors.New("require project name")
		}
		opts := todoist.NewProjectOpts{}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			opts.Color = color
		}
		if parentStr, err := cmd.Flags().GetString("parent"); err != nil {
			return err
		} else {
			if len(parentStr) != 0 {
				if parent, err := todoist.NewID(parentStr); err != nil {
					return err
				} else {
					opts.ParentID = parent
				}
			}
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			return err
		} else {
			opts.ChildOrder = order
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			return err
		} else {
			opts.IsFavorite = todoist.IntBool(favorite)
		}
		project, err := todoist.NewProject(name, &opts)
		if project == nil {
			return errors.New("failed to initialize a project")
		}
		if err != nil {
			return err
		}
		if _, err = client.Project.Add(*project); err != nil {
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
			return errors.New("failed to add this project. it may be failed to sync")
		}
		// it may not be new project
		syncedProject := projects[len(projects)-1]
		fmt.Println("succeeded to add a project")
		fmt.Println(util.ProjectTableString([]todoist.Project{syncedProject}))
		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "update project",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := util.NewClient()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			return errors.New("require project id to update")
		}
		id, err := todoist.NewID(args[0])
		if err != nil {
			return fmt.Errorf("invalid id: %s", args[0])
		}
		project := client.Project.Resolve(id)
		if project == nil {
			return fmt.Errorf("no such project id: %s", id)
		}
		if name, err := cmd.Flags().GetString("name"); err != nil {
			return err
		} else {
			if len(name) != 0 {
				project.Name = name
			}
		}
		if color, err := cmd.Flags().GetInt("color"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("color") {
				project.Color = color
			}
		}
		if collapsed, err := cmd.Flags().GetBool("collapsed"); err != nil {
			return err
		} else {
			if collapsed {
				project.Collapsed = true
			}
		}
		if unCollapsed, err := cmd.Flags().GetBool("un-collapsed"); err != nil {
			return err
		} else {
			if unCollapsed {
				project.Collapsed = false
			}
		}
		if favorite, err := cmd.Flags().GetBool("favorite"); err != nil {
			return err
		} else {
			if favorite {
				project.IsFavorite = true
			}
		}
		if unFavorite, err := cmd.Flags().GetBool("un-favorite"); err != nil {
			return err
		} else {
			if unFavorite {
				project.IsFavorite = false
			}
		}
		if _, err = client.Project.Update(*project); err != nil {
			return err
		}
		if parentStr, err := cmd.Flags().GetString("parent"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("parent") {
				if parent, err := todoist.NewID(parentStr); err != nil {
					return err
				} else {
					if err = client.Project.Move(project.ID, parent); err != nil {
						return err
					}
				}
			}
		}
		if order, err := cmd.Flags().GetInt("order"); err != nil {
			return err
		} else {
			if cmd.Flags().Changed("order") {
				project.ChildOrder = order
				if err = client.Project.Reorder([]todoist.Project{*project}); err != nil {
					return err
				}
			}
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
		fmt.Println("succeeded to update the project")
		fmt.Println(util.ProjectTableString([]todoist.Project{*syncedProject}))
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			if len(args) == 0 {
				return errors.New("require project id to delete")
			}
			return util.ProcessID(args[0], func(id todoist.ID) error {
				project := client.Project.Resolve(id)
				if project == nil {
					return fmt.Errorf("invalid project id: %s", id)
				}
				fmt.Println(util.ProjectTableString([]todoist.Project{*project}))
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("are you sure to delete above project? (y/[n]): ")
				ans, err := reader.ReadString('\n')
				if ans != "y\n" || err != nil {
					fmt.Println("abort")
					return errors.New("abort")
				}
				return client.Project.Delete(id)
			})
		}); err != nil {
			if err.Error() == "abort" {
				return nil
			}
			return err
		}
		fmt.Println("succeeded to delete the project")
		return nil
	},
}

var projectArchiveCmd = &cobra.Command{
	Use:   "archive [id]",
	Short: "archive project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			if len(args) == 0 {
				return errors.New("require project id to archive")
			}
			return util.ProcessID(args[0], func(id todoist.ID) error {
				return client.Project.Archive(id)
			})
		}); err != nil {
			return err
		}
		fmt.Println("succeeded to archive the project")
		return nil
	},
}

var projectUnarchiveCmd = &cobra.Command{
	Use:   "unarchive [id]",
	Short: "unarchive project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := util.AutoCommit(func(client todoist.Client, ctx context.Context) error {
			if len(args) == 0 {
				return errors.New("require project id to un-archive")
			}
			return util.ProcessID(args[0], func(id todoist.ID) error {
				return client.Project.Unarchive(id)
			})
		}); err != nil {
			return err
		}
		fmt.Println("succeeded to un-archive the project")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectListCmd)
	projectAddCmd.Flags().IntP("color", "c", 47, "color")
	projectAddCmd.Flags().String("parent", "", "parent project id")
	projectAddCmd.Flag("parent").Annotations = map[string][]string{cobra.BashCompCustom: {"__todoist_project_id"}}
	projectAddCmd.Flags().Int("order", 0, "child order")
	projectAddCmd.Flags().Bool("favorite", false, "is favorite")
	projectCmd.AddCommand(projectAddCmd)
	projectUpdateCmd.Flags().String("name", "", "name of the project")
	projectUpdateCmd.Flags().IntP("color", "c", 47, "color")
	projectUpdateCmd.Flags().String("parent", "", "parent project id")
	projectUpdateCmd.Flag("parent").Annotations = map[string][]string{cobra.BashCompCustom: {"__todoist_project_id"}}
	projectUpdateCmd.Flags().Int("order", 0, "child order")
	projectUpdateCmd.Flags().Bool("collapsed", false, "collapse project")
	projectUpdateCmd.Flags().Bool("un-collapsed", false, "un-collapse project")
	projectUpdateCmd.Flags().Bool("favorite", false, "is favorite")
	projectUpdateCmd.Flags().Bool("un-favorite", false, "is not favorite")
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectUnarchiveCmd)
}
