package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/zyedidia/generic/multimap"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}

type valueSource interface {
	orig() string
	values() []string
}

type model struct {
	orig        string
	values      []string
	pathsLookup multimap.MultiMap[string, string]
	paginator   paginator.Model
}

type envSource struct {
	name string
}

func newEnvSource() *envSource {
	if len(os.Args) != 2 {
		fmt.Println("please provide the name of the environment variable to debug")
		os.Exit(1)
	}
	return &envSource{os.Args[1]}
}

func (s *envSource) values() []string {
	return strings.Split(s.orig(), string(os.PathListSeparator))
}

func (s *envSource) orig() string {
	return os.Getenv(s.name)
}

func initialModel() model {

	var source valueSource = newEnvSource()
	values := source.values()

	pathsLookup := multimap.NewMapSet[string](func(a, b string) bool { return true })
	for _, path := range values {
		pathsLookup.Put(path, path)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(values))

	return model{
		orig:        source.orig(),
		values:      values,
		pathsLookup: pathsLookup,
		paginator:   p,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		const margin = 6
		m.paginator.PerPage = int(math.Max(1, float64(msg.Height-margin)))
		// re-calculate the paging state
		m.paginator.SetTotalPages(len(m.values))
		// try to keep the page but choose another if necessary
		m.paginator.Page = int(math.Min(math.Max(0, float64(m.paginator.Page)), float64(m.paginator.TotalPages-1)))

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(`tap Esc/q/Ctrl-C to quit
`)

	start, end := m.paginator.GetSliceBounds(len(m.values))

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Dup #", "Bad", "Path"})
	for _, path := range m.values[start:end] {
		rep := " "
		count := len(m.pathsLookup.Get(path))
		if count > 1 {
			rep = fmt.Sprint(count)
		}
		t.AppendRow(table.Row{
			rep,
			statusOf(path),
			path,
		})
	}
	b.WriteString(t.Render())
	b.WriteString("\n")

	b.WriteString("  " + m.paginator.View())

	return b.String()
}

func statusOf(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "X"
	}

	if !fileInfo.IsDir() {
		return "F"
	}

	return " "
}
