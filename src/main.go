package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	psutil "github.com/shirou/gopsutil/v4/process"
)

type sortKey int

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

type process struct {
	Pid  int
	Name string
	Args string
	RSS  uint64
	CPU  float64
}

type refreshMsg []process

type model struct {
	table   table.Model
	sortKey sortKey
	height  int
	width   int
}

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
		{Title: "CPU", Width: 5},
		{Title: "Command", Width: 50}, // resizable dynamically
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
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return refreshMsg(getProcesses())
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
			{Title: "Mem (KB)", Width: 8},
			{Title: "CPU", Width: 3},
			{Title: "Command", Width: cmdWidth},
		})

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "N":
			m.sortKey = SortKeyName
			m.updateTable(getProcesses())
		case "M":
			m.sortKey = SortKeyMemory
			m.updateTable(getProcesses())
		case "P":
			m.sortKey = SortKeyCPU
			m.updateTable(getProcesses())
		case "cltr+p", "up", "k":
			m.table.MoveUp(1)
		case "ctrl+n", "down", "j":
			m.table.MoveDown(1)
		}

	case refreshMsg:
		m.updateTable([]process(msg))
		return m, refreshCmd()
	}

	return m, nil
}

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	line := headerStyle.Render(
		fmt.Sprintf(
			"ytop — sorting by %s | q to quit",
			m.sortKey.String(),
		))

	return line + "\n\n" + m.table.View()
}

func getProcesses() []process {
	procs, err := psutil.Processes()
	if err != nil {
		return nil
	}

	var result []process

	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		args, _ := p.Cmdline()

		mem, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		cpu, err := p.CPUPercent()
		if err != nil {
			continue
		}

		cmd := args
		if cmd == "" {
			cmd = name
		}

		result = append(result, process{
			Pid:  int(p.Pid),
			Name: name,
			Args: cmd,
			RSS:  mem.RSS / 1024,
			CPU:  cpu,
		})
	}

	return result
}

func (m *model) updateTable(procs []process) {
	m.sortProcesses(procs)

	rows := make([]table.Row, len(procs))
	for i, p := range procs {
		rows[i] = table.Row{
			fmt.Sprintf("%d", p.Pid),
			fmt.Sprintf("%d", p.RSS),
			fmt.Sprintf("%.1f", p.CPU),
			p.Args,
		}
	}
	m.table.SetRows(rows)
}

func (m *model) sortProcesses(processes []process) {
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

func main() {
	m := model{
		table:   newProcessTable(),
		sortKey: SortKeyCPU,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
