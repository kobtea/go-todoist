package cmd

import (
	"fmt"
	"github.com/kobtea/go-todoist/cmd/util"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/spf13/cobra"
	"sort"
	"strconv"
	"strings"
)

// todayCmd represents the today command
var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "show today's tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}
		items := client.Item.FindByDueDate(todoist.Today())
		sort.Slice(items, func(i, j int) bool {
			return items[i].DueDateUtc.Before(items[j].DueDateUtc)
		})
		relations := client.Relation.Items(items)
		var rows [][]todoist.ColorStringer
		for _, i := range items {
			project := todoist.Project{}
			if v, ok := relations.Projects[i.ProjectID]; ok {
				project = v
			}
			labels := []string{}
			for _, lid := range i.Labels {
				if v, ok := relations.Labels[lid]; ok {
					labels = append(labels, v.String())
				}
			}
			rows = append(rows, []todoist.ColorStringer{
				todoist.NewNoColorString(i.ID.String()),
				i.DueDateUtc,
				todoist.NewNoColorString(strconv.Itoa(i.Priority)),
				project,
				todoist.NewNoColorString(strings.Join(labels, " ")),
				todoist.NewNoColorString(i.Content),
			})
		}
		fmt.Println(util.TableString(rows))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(todayCmd)
}
