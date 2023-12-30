package view

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/d-led/pathdebug/render"
)

func RenderCSV() string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.Write([]string{"Id", "Duplicates", "Bad", "Path"})

	for _, row := range getResults() {
		w.Write([]string{
			fmt.Sprint(row.Id),
			render.FormatList(row.Duplicates),
			render.StatusOfDir(row),
			row.Path,
		})
	}
	w.Flush()
	return buf.String()
}
