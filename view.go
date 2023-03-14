package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-pretty/v6/table"
)

type viewModel struct {
	results   []resultRow
	paginator paginator.Model
}

func initialModel() viewModel {
	if len(os.Args) != 2 {
		failWith("please provide the name of the environment variable to debug")
	}
	source := newEnvSource(os.Args[1])

	fs := &osFilesystem{}

	results, err := calculateResults(fs, source)
	if err != nil {
		failWith(err.Error())
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.SetTotalPages(len(results))

	return viewModel{
		results:   results,
		paginator: p,
	}
}

func (m viewModel) Init() tea.Cmd {
	return nil
}

func (m viewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		const margin = 6
		m.paginator.PerPage = int(math.Max(1, float64(msg.Height-margin)))
		// re-calculate the paging state
		m.paginator.SetTotalPages(len(m.results))
		// try to keep the page but choose another if necessary
		m.paginator.Page = int(math.Min(math.Max(0, float64(m.paginator.Page)), float64(m.paginator.TotalPages-1)))

	case tea.KeyMsg:

		switch msg.String() {

		// quit
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)

	return m, cmd
}

func (m viewModel) View() string {
	var b strings.Builder

	b.WriteString(`tap Esc/q/Ctrl-C to quit, <-/-> to paginate
`)

	start, end := m.paginator.GetSliceBounds(len(m.results))
	results := m.results[start:end]

	m.renderTable(&b, results, start)

	return b.String()
}

func (m viewModel) renderTable(b *strings.Builder, results []resultRow, offset int) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Dup[#]", "Bad", "Path"})
	for _, row := range results {
		t.AppendRow(table.Row{
			row.id,
			formatDuplicates(row.duplicates),
			statusOfDir(row),
			row.path,
		})
	}
	b.WriteString(t.Render())
	b.WriteString("\n")

	b.WriteString("  " + m.paginator.View())
}

func statusOfDir(row resultRow) string {
	if !row.exists {
		return "X"
	}

	if !row.isDir {
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
