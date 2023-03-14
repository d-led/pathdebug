package view

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/d-led/pathdebug/common"
	"github.com/d-led/pathdebug/render"
)

func Run() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

type viewModel struct {
	results   []common.ResultRow
	paginator paginator.Model
}

func initialModel() viewModel {
	results := getResults()

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

	render.RenderTable(&b, results)

	b.WriteString("\n")
	b.WriteString("  " + m.paginator.View())

	return b.String()
}
