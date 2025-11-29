package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
	"skybert.net/ytop/pkg"
)

// Populated at build time
var Version = "dev"

type (
	refreshMsg []pkg.Process
	sortKey    int
	mode       int
)

const (
	headerHeight          = 5
	updateIntervalSeconds = 2

	modeViewTable mode = iota
	modeViewProcess
	modeSearchProcess
)

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
	conf       pkg.YTopConf
	height     int
	humanSizes bool
	mode       mode
	processes  []pkg.Process
	simpleView bool
	sortKey    pkg.SortKey
	table      table.Model
	width      int
}

func (m model) Init() tea.Cmd {
	return m.refreshCmd()
}

func (m model) refreshCmd() tea.Cmd {
	return tea.Tick(
		time.Second*time.Duration(updateIntervalSeconds),
		func(t time.Time) tea.Msg {
			m.processes = Processes()
			return refreshMsg(m.processes)
		})
}

func (m *model) updateTable(procs []pkg.Process) {
	SortProcesses(procs, m.sortKey)

	rows := make([]table.Row, len(procs))
	for i, p := range procs {
		row := table.Row{
			fmt.Sprintf("%d", p.Pid),
			m.humanBytes(p.RSS),
			fmt.Sprintf("%.1f", p.CPU),
			p.Name,
		}
		if !m.simpleView {
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
		m.width = msg.Width
		m.height = msg.Height
		m.table = m.createTable()
		m.table.SetHeight(max(msg.Height-headerHeight-1, 5))
		columns := TableColumns(m.simpleView, msg.Width)
		m.table.SetColumns(columns)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			m.humanSizes = !m.humanSizes
			m.updateTable(Processes())
		case "N", "n":
			m.sortKey = pkg.SortKeyName
			m.updateTable(Processes())
		case "M", "m":
			m.sortKey = pkg.SortKeyMemory
			m.updateTable(Processes())
		case "P", "p":
			m.sortKey = pkg.SortKeyCPU
			m.updateTable(Processes())
		case "S", "s":
			m.simpleView = !m.simpleView
			m.table = m.createTable()
			m.table.SetHeight(m.height - headerHeight)
			m.updateTable(Processes())
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+p", "up", "k":
			m.table.MoveUp(1)
		case "ctrl+n", "down", "j":
			m.table.MoveDown(1)
		}
	case refreshMsg:
		m.updateTable([]pkg.Process(msg))
		return m, m.refreshCmd()
	}

	return m, nil
}

func (m model) View() string {
	return ViewHeader(m.sortKey, m.humanSizes) +
		"\n" +
		m.table.View()
}

func main() {
	pflag.Parse()
	if showVersion {
		fmt.Printf("ytop version: %v\n", Version)
		os.Exit(0)
	}

	m := model{
		conf: pkg.YTopConf{
			HeaderForeground:   "",
			HeaderBackground:   "",
			SelectedForeground: fgSelColour,
			SelectedBackground: bgSelColour,
			SimpleView:         simpleView,
		},
		table:      table.Model{},
		humanSizes: humanSizes,
		simpleView: simpleView,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("%v: %v\n", "There was an error", err)
		fmt.Printf("%v\n", debug.Stack())
		os.Exit(1)
	}

}
