package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"skybert.net/ytop/pkg"
)

func (m *model) viewHeader() string {
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
			HumanUptime(),
			m.sortKey.String(),
		))
	cpu := helpStyle.Render("CPU: ") + CPUSummary()
	memory := helpStyle.Render("Memory: ") + MemorySummary(m.humanSizes)
	help := helpStyle.Render(
		fmt.Sprintf(
			"Sort: %s %s %s | %s ",
			"p (cpu)",
			"m (mem)",
			"n (name)",
			"quit: q",
		))

	return line + "\n" + cpu + "\n" + memory + "\n" + help + "\n"
}

func tableHeaders(simpleView bool) []string {
	columns := []string{
		"PID",
		"RSS",
		"%CPU",
	}

	if simpleView {
		columns = append(columns, "NAME")
	} else {
		columns = append(columns, "NAME")
		columns = append(columns, "COMMAND")
	}

	return columns
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

	procsToInclude := make([]pkg.Process, 0)
	for _, p := range procs {
		if m.searchQuery != "" {
			if !strings.Contains(
				strings.ToLower(p.Args),
				strings.ToLower(m.searchQuery)) {
				continue
			}
		}
		procsToInclude = append(procsToInclude, p)
	}

	rows := make([][]string, len(procsToInclude))
	for i, p := range procsToInclude {
		row := []string{
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

	m.table.ClearRows().Rows(rows...)
}

func (m *model) humanBytes(bytes uint64) string {
	if !m.humanSizes {
		// Default is showing size in bytes
		return fmt.Sprintf("%v", bytes)
	}

	return pkg.HumanBytes(bytes)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table = m.createTable()
		m.table.Height(max(msg.Height-headerHeight-1, 5))
		headers := tableHeaders(m.simpleView)
		m.table.Headers(headers...)
		return m, nil

	case tea.KeyPressMsg:
		if m.searchShow {
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			cmds = append(cmds, cmd)

			switch msg.String() {
			case "enter":
				m.searchInput.SetValue("")
				m.searchShow = false
				log.Println("You searched for: " + m.searchQuery)
			case "esc":
				m.searchShow = false
				m.searchInput.SetValue("")
				m.searchQuery = ""
				m.updateTable(Processes())
			default:
				m.searchQuery = m.searchInput.Value()
				m.updateTable(Processes())
			}

			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "/":
			m.searchQuery = ""
			m.searchShow = true
			cmd := m.searchInput.Focus()
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)

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
			m.table.Headers(tableHeaders(m.simpleView)...)
			m.updateTable(Processes())
		case "esc":
			// Escape cancels any active search filter
			m.searchQuery = ""
			m.updateTable(Processes())
		case "ctrl+c", "q":
			return m, tea.Quit
			// case "ctrl+p", "up", "k":
			// 	m.table.MoveUp(1)
			// case "ctrl+n", "down", "j":
			// 	m.table.MoveDown(1)
		}
	case refreshMsg:
		// There are two things going on here on the update
		// loop: We want to update the table as well as update
		// the input field where the user is typing, hence the
		// cmds array.
		cmds = append(cmds, m.refreshCmd())
		m.updateTable(Processes())

		if m.searchShow {
			cmd := m.searchInput.Focus()
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() tea.View {
	// Regular top table
	s := lipgloss.JoinVertical(
		lipgloss.Top,
		m.viewHeader()+"\n"+m.table.String()+"\n")
	// Showing the search input box
	if m.searchShow {
		m.searchInput.Focus()
		s = lipgloss.JoinVertical(
			lipgloss.Top,
			m.viewHeader()+
				m.searchInput.View()+"\n"+
				m.table.String()+"\n")
	}

	v := tea.NewView(s)

	var c *tea.Cursor
	if !m.searchInput.VirtualCursor() {
		c = m.searchInput.Cursor()
	}
	v.Cursor = c

	// Alternate screen buffer (AltScreen) means full screen 😉
	v.AltScreen = true
	return v
}

func (m *model) createTable() *table.Table {
	headers := tableHeaders(m.simpleView)
	dummyRow := make([]string, len(headers))
	dummyRow[0] = "Loading processes ..."

	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.
		Background(lipgloss.Color(m.conf.HeaderBackground)).
		Foreground(lipgloss.Color(m.conf.HeaderForeground)).
		Bold(true)

	t := table.New().
		Headers(headers...).
		Rows(dummyRow).
		BorderBottom(false).
		BorderColumn(false).
		BorderHeader(false).
		BorderLeft(false).
		BorderRight(false).
		BorderTop(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return headerStyle
			default:
				return baseStyle.
					Foreground(lipgloss.Color(m.conf.Foreground))
			}
		})

	return t
}

func searchInput() textinput.Model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Cmd to search for"
	searchInput.CharLimit = 100
	searchInput.SetWidth(searchInput.CharLimit)
	searchInput.Focus()
	return searchInput
}
