package cmd

import (
	"github.com/fatih/color"
	"github.com/kobtea/go-todoist/todoist"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
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
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"id", "date", "pri", "project", "labels", "content"})
		table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
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
			date := i.DateString
			if i.IsOverDueDate() {
				date = color.New(color.BgRed).SprintFunc()(date)
			}
			table.Append([]string{
				i.ID.String(),
				date,
				strconv.Itoa(i.Priority),
				project,
				strings.Join(labels, " "),
				i.Content,
			})

		}
		table.Render()
		return nil
	},
}

func init() {
	RootCmd.AddCommand(todayCmd)
}
