package util

import (
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
)

func StringWidthWithoutColor(s string) int {
	re := regexp.MustCompile("\x1b\\[[0-9;]*m")
	return runewidth.StringWidth(re.ReplaceAllString(s, ""))
}

func TableString(rows [][]string) string {
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
			// FIXME: regex may be slow
			l := StringWidthWithoutColor(ss[i])
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
			if _, err := strconv.Atoi(rows[i][j]); err == nil {
				f = runewidth.FillLeft
			}
			res += f(rows[i][j], lens[j])
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
