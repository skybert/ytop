package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"skybert.net/ytop/internal"
	"skybert.net/ytop/pkg"
)

type mode int
type sortKey int

const (
	modeViewTable mode = iota
	modeViewProcess
)

const (
	SortKeyCPU sortKey = iota
	SortKeyMemory
	SortKeyName
)

func (k sortKey) String() string {
	switch k {
	case SortKeyCPU:
		return "cpu"
	case SortKeyName:
		return "name"
	case SortKeyMemory:
		return "memory"
	}
	return "unknown"
}

type refreshMsg []pkg.Process

type model struct {
	table   table.Model
	sortKey sortKey
	height  int
	width   int
	mode    mode
}

var updateIntervalSeconds = 2

func newStyles() table.Styles {
	s := table.DefaultStyles()

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#222235")).
		Background(lipgloss.Color("#06c993")).
		Bold(true)

	return s
}

func newProcessTable() table.Model {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "Mem (KB)", Width: 10},
		{Title: "CPU", Width: 7},
		{Title: "Command", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	t.SetStyles(newStyles())
	return t
}

func (m model) Init() tea.Cmd {
	return refreshCmd()
}

func refreshCmd() tea.Cmd {
	return tea.Tick(
		time.Second*time.Duration(updateIntervalSeconds),
		func(t time.Time) tea.Msg {
			return refreshMsg(internal.Processes())
		})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3 // top bar height
		tableHeight := max(m.height-headerHeight-1, 5)

		m.table.SetHeight(tableHeight)

		// Expand Command column to fill available width
		cmdWidth := max(m.width-30, 20)
		m.table.SetColumns([]table.Column{
			{Title: "PID", Width: 7},
			{Title: "RES", Width: 8},
			{Title: "%CPU", Width: 5},
			{Title: "NAME", Width: 20},
			{Title: "COMMAND", Width: cmdWidth},
		})

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selectedRowContainsData := len(m.table.SelectedRow()) > 0
			if m.mode == modeViewTable && selectedRowContainsData {
				m.mode = modeViewProcess
			}
		case "esc":
			if m.mode == modeViewProcess {
				m.mode = modeViewTable
			}
		case "N", "n":
			m.sortKey = SortKeyName
			m.updateTable(internal.Processes())
		case "M", "m":
			m.sortKey = SortKeyMemory
			m.updateTable(internal.Processes())
		case "P", "p":
			m.sortKey = SortKeyCPU
			m.updateTable(internal.Processes())
		case "ctrl+p", "up", "k":
			m.table.MoveUp(1)
		case "ctrl+n", "down", "j":
			m.table.MoveDown(1)
		}

	case refreshMsg:
		m.updateTable([]pkg.Process(msg))
		return m, refreshCmd()
	}

	return m, nil
}

func (m model) viewHeader() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)
	helpStyle := lipgloss.NewStyle().
		Bold(false).
		Padding(0, 1)

	line := headerStyle.Render(
		fmt.Sprintf(
			"ytop — %s up %s | sorted by %s",
			time.Now().Format("15:04:05"),
			internal.HumanUptime(),
			m.sortKey.String(),
		))
	help := helpStyle.Render(
		fmt.Sprintf(
			"Sort: %s %s %s | %s ",
			"p (cpu)",
			"m (mem)",
			"n (name)",
			"quit: q",
		))
	return line + "\n" + help + "\n\n"

}

func (m model) viewTable() string {
	return m.viewHeader() + m.table.View()

}

func (m model) viewProcess() string {
	row := m.table.SelectedRow()
	if len(row) == 0 {
		return ""
	}
	pid, err := strconv.Atoi(row[0])
	if err != nil {
		return ""
	}
	p := internal.Process(pid)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	sort.Slice(p.Env, func(i, j int) bool {
		one := p.Env[i]
		two := p.Env[j]
		return one < two
	})

	procInfo := labelStyle.Render("PID:") + " " + strconv.Itoa(p.Pid) + "\n" +
		labelStyle.Render("Name:") + " " + p.Name + "\n" +
		labelStyle.Render("Command with arguments:") + " " + strings.ReplaceAll(p.Args, " ", "\n   ") + "\n" +
		labelStyle.Render("Unix env vars:") + " " + strings.Join(p.Env, "\n  ")

	return m.viewHeader() + procInfo + "\n"
}

func (m model) View() string {
	switch m.mode {
	case modeViewTable:
		return m.viewTable()
	case modeViewProcess:
		return m.viewProcess()
	}
	return ""
}

func (m *model) sortProcesses(processes []pkg.Process) {
	sort.Slice(processes, func(i, j int) bool {
		pi := processes[i]
		pj := processes[j]

		switch m.sortKey {
		case SortKeyMemory:
			return pi.RSS > pj.RSS
		case SortKeyCPU:
			return pi.CPU > pj.CPU
		case SortKeyName:
			// Sort name ascending
			return pi.Name < pj.Name
		}
		return false
	})
}

func (m *model) updateTable(procs []pkg.Process) {
	m.sortProcesses(procs)

	rows := make([]table.Row, len(procs))
	for i, p := range procs {
		rows[i] = table.Row{
			fmt.Sprintf("%d", p.Pid),
			fmt.Sprintf("%d", p.RSS),
			fmt.Sprintf("%.1f", p.CPU),
			p.Name,
			p.Args,
		}
	}
	m.table.SetRows(rows)
}

func main() {
	m := model{
		table:   newProcessTable(),
		sortKey: SortKeyCPU,
		mode:    modeViewTable,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
