package render

import (
	"fmt"
	"strings"

	"github.com/d-led/pathdebug/common"
	"github.com/jedib0t/go-pretty/v6/table"
)

func RenderTableToString(results []common.ResultRow) string {
	var b strings.Builder
	RenderTable(&b, results)
	return b.String()
}

func RenderTable(b *strings.Builder, results []common.ResultRow) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Dup[#]", "Bad", "Path"})
	for _, row := range results {
		t.AppendRow(table.Row{
			row.Id,
			formatDuplicates(row.Duplicates),
			statusOfDir(row),
			row.Path,
		})
	}
	b.WriteString(t.Render())
}

func statusOfDir(row common.ResultRow) string {
	if !row.Exists {
		return "X"
	}

	if !row.IsDir {
		return "F"
	}

	return " "
}

func formatDuplicates(ids []int) string {
	res := []string{}
	for _, id := range ids {
		res = append(res, fmt.Sprint(id))
	}
	return strings.Join(res, ", ")
}
