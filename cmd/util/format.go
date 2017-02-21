package util

import (
	"github.com/kobtea/go-todoist/todoist"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

func StringWidthWithoutColor(s string) int {
	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return runewidth.StringWidth(re.ReplaceAllString(s, ""))
}

func TableString(rows [][]todoist.ColorStringer) string {
	if len(rows) == 0 {
		return ""
	}
	min := len(rows[0])
	for _, i := range rows {
		if len(i) < min {
			min = len(i)
		}
	}
	// retrieve max width in each columns
	lens := make([]int, min)
	for _, ss := range rows {
		for i := 0; i < len(lens); i++ {
			l := runewidth.StringWidth(ss[i].String())
			if l > lens[i] {
				lens[i] = l
			}
		}
	}
	// format rows into table
	var res string
	for i := 0; i < len(rows); i++ {
		for j := 0; j < len(lens); j++ {
			f := runewidth.FillRight
			if _, err := strconv.Atoi(rows[i][j].String()); err == nil {
				f = runewidth.FillLeft
			}
			colorSeqLen := runewidth.StringWidth(rows[i][j].ColorString()) - runewidth.StringWidth(rows[i][j].String())
			res += f(rows[i][j].ColorString(), lens[j]+colorSeqLen)
			if j < len(lens)-1 {
				res += " "
			}
		}
		if i < len(rows)-1 {
			res += "\n"
		}
	}
	return res
}

func ItemTableString(items []todoist.Item, relations todoist.ItemRelations) string {
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
	return TableString(rows)
}