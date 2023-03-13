package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mitchellh/go-homedir"
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
	source() string
	orig() string
	values() []string
}

type model struct {
	values      []string
	pathsLookup multimap.MultiMap[string, string]
	paginator   paginator.Model
}

type envSource struct {
	name string
}

func newEnvSource() *envSource {
	if len(os.Args) != 2 {
		failWith("please provide the name of the environment variable to debug")
	}
	return &envSource{os.Args[1]}
}

func (s *envSource) values() []string {
	paths := strings.Split(s.orig(), string(os.PathListSeparator))
	if len(paths) == 1 && paths[0] == "" {
		return []string{}
	}
	return paths
}

func (s *envSource) orig() string {
	return os.Getenv(s.name)
}

func (s *envSource) source() string {
	return s.name
}

func initialModel() model {
	var source valueSource = newEnvSource()
	values := source.values()

	if len(values) == 0 {
		failWith(source.source() + " is empty")
	}

	duplicatePredicate := func(a, b string) bool { return true }
	pathsLookup := multimap.NewMapSet[string](duplicatePredicate)
	for _, path := range values {
		pathKey := getAbsolutePath(path)
		pathsLookup.Put(pathKey, path)
	}

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.SetTotalPages(len(values))

	return model{
		values:      values,
		pathsLookup: pathsLookup,
		paginator:   p,
	}
}

func getAbsolutePath(path string) string {
	expandedPath, err := homedir.Expand(path)
	if err == nil {
		path = expandedPath
	}
	absPath, err := filepath.Abs(path)
	if err == nil {
		path = absPath
	}
	return os.ExpandEnv(path)
}

func (m model) Init() tea.Cmd {
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

		switch msg.String() {

		// quit
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(`tap Esc/q/Ctrl-C to quit, <-/-> to paginate
`)

	start, end := m.paginator.GetSliceBounds(len(m.values))
	paths := m.values[start:end]

	m.renderTable(&b, paths)

	return b.String()
}

func (m model) renderTable(b *strings.Builder, paths []string) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Dup #", "Bad", "Path"})
	for _, path := range paths {
		rep := " "
		count := len(m.pathsLookup.Get(getAbsolutePath(path)))
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
}

func statusOf(path string) string {
	path = getAbsolutePath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "X"
	}

	if !fileInfo.IsDir() {
		return "F"
	}

	return " "
}

func failWith(message string) {
	fmt.Println(message)
	os.Exit(1)
}
