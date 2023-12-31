package render

import (
	"fmt"
	"runtime"
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
	headers := table.Row{"#", "Dup[#]", "Bad", "Path"}
	addSources := runtime.GOOS != "windows"
	if addSources {
		headers = append(headers, "±Sources")
	}
	t.AppendHeader(headers)
	for _, row := range results {
		r := table.Row{
			row.Id,
			FormatList(row.Duplicates),
			StatusOfDir(row),
			row.Path,
		}
		if addSources {
			r = append(r, FormatList(row.CandidateSources))
		}
		t.AppendRow(r)
	}
	b.WriteString(t.Render())
}

func StatusOfDir(row common.ResultRow) string {
	if !row.Exists {
		return "X"
	}

	if !row.IsDir {
		return "F"
	}

	return " "
}

func FormatList[T any](ids []T) string {
	res := []string{}
	for _, id := range ids {
		res = append(res, fmt.Sprint(id))
	}
	return strings.Join(res, ", ")
}
