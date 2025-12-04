package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
	"skybert.net/ytop/internal"
	"skybert.net/ytop/pkg"
)

// Populated at build time
var Version = "dev"

var humanSizes bool
var showVersion bool
var simpleView bool
var bgSelColour string
var fgSelColour string

func init() {
	pflag.BoolVarP(
		&humanSizes,
		"human-readable",
		"h",
		false,
		"Human readable sizes in chunks of 1024")
	pflag.BoolVarP(&showVersion, "version", "v", false, "Show version")
	pflag.BoolVarP(&simpleView, "simple", "s", false, "Simple view, less info")
	pflag.StringVar(&fgSelColour, "sel-fg", "#222235", "Selection background colour")
	pflag.StringVar(&bgSelColour, "sel-bg", "#06c993", "Selection foreground colour")
}

type model struct {
	table      table.Model
	processes  []pkg.Process
	humanSizes bool
}

func (m model) Init() tea.Cmd {
	return m.refreshCmd()
}

func (m model) refreshCmd() tea.Cmd {
	return tea.Tick(
		time.Second*time.Duration(updateIntervalSeconds),
		func(t time.Time) tea.Msg {
			m.processes = internal.Processes()
			return refreshMsg(m.processes)
		})
}

const (
	headerHeight          = 4
	updateIntervalSeconds = 2
)

type refreshMsg []pkg.Process

func (m *model) updateTable(procs []pkg.Process) {
	rows := make([]table.Row, len(procs))
	for i, p := range procs {
		row := table.Row{
			fmt.Sprintf("%d", p.Pid),
			m.humanBytes(p.RSS),
			fmt.Sprintf("%.1f", p.CPU),
			p.Name,
		}
		if !simpleView {
			row = append(row, p.Args)
		}
		rows[i] = row
	}

	m.table.SetRows(rows)
}

func (m *model) humanBytes(bytes uint64) string {
	if !m.humanSizes {
		// Default is showing size in bytes
		return fmt.Sprintf("%v", bytes)
	}

	return pkg.HumanBytes(bytes)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)
		m.table.SetHeight(max(msg.Height-headerHeight-1, 5))
		columns := internal.TableColumns(simpleView, msg.Width)
		m.table.SetColumns(columns)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case refreshMsg:
		m.updateTable([]pkg.Process(msg))
		return m, m.refreshCmd()
	}

	return m, nil
}

func (m model) View() string {
	header := "ytop"
	info := "foo bar baz info"
	return header + "\n" + info + "\n\n" + m.table.View()
}

func main() {
	pflag.Parse()
	if showVersion {
		fmt.Printf("ytop version: %v\n", Version)
		os.Exit(0)
	}

	m := model{
		table: internal.CreateTable(simpleView, 80),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("%v: %v\n", "There was an error", err)
		fmt.Printf("%v\n", debug.Stack())
		os.Exit(1)
	}

}
