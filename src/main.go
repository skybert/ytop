package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shirou/gopsutil/v4/process"
)

type sortKey int

const (
	SortKeyCPU sortKey = iota
	SortKeyMemory
)

func (k sortKey) String() string {
	switch k {
	case SortKeyCPU:
		return "cpu"
	case SortKeyMemory:
		return "memory"
	default:
		return "unknown"
	}
}

type model struct {
	processes []*process.Process
	sortKey   sortKey
}

func initialModel() model {
	procs, err := process.Processes()
	if err != nil {
		panic(err)
	}

	return model{
		processes: procs,
	}

}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		case "M":
			m.sortKey = SortKeyMemory
		case "P":
			m.sortKey = SortKeyCPU
		}
	}

	return m, nil
}

func (m model) View() string {
	maxProc := 20
	s := fmt.Sprintf(
		"ytop showing top %d processes sorted by %s\n",
		maxProc,
		m.sortKey.String(),
	)

	sort.Slice(m.processes, func(i, j int) bool {
		pi := m.processes[i]
		pj := m.processes[j]

		switch m.sortKey {
		case SortKeyMemory:
			mi, err1 := pi.MemoryInfo()
			mj, err2 := pj.MemoryInfo()
			if err1 != nil && err2 == nil {
				return false
			}
			if err2 != nil && err1 == nil {
				return true
			}
			if err1 != nil && err2 != nil {
				return false // maintain relative order
			}
			return mi.RSS > mj.RSS
		case SortKeyCPU:
			ci, err1 := pi.CPUPercent()
			cj, err2 := pj.CPUPercent()
			if err1 != nil && err2 == nil {
				return false
			}
			if err2 != nil && err1 == nil {
				return true
			}
			if err1 != nil && err2 != nil {
				return false // maintain relative order
			}
			return ci > cj
		}
		return false
	})

	for i, proc := range m.processes {
		if i > maxProc {
			break
		}

		name, _ := proc.Cmdline()
		mem, _ := proc.MemoryInfo()
		memPercent, _ := proc.MemoryPercent()

		s += fmt.Sprintf("%d", proc.Pid) +
			" " + strconv.FormatUint(mem.RSS, 10) +
			" " + fmt.Sprintf("%.1f", memPercent) +
			" " + name +
			"\n"
	}

	s += "\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
