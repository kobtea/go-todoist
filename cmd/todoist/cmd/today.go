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
		var rows [][]string
		for _, i := range items {
			var project string
			if p := client.Project.Resolve(i.ProjectID); p != nil {
				project = p.String()
			}
			var labels []string
			for _, l := range i.Labels {
				if j := client.Label.Resolve(l); j != nil {
					labels = append(labels, j.String())
				}
			}
			rows = append(rows, []string{
				i.ID.String(),
				i.DueDateUtc.Local().ColorShortString(),
				strconv.Itoa(i.Priority),
				project,
				strings.Join(labels, " "),
				i.Content,
			})
		}
		fmt.Println(util.TableString(rows))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(todayCmd)
}
