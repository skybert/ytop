package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/spf13/pflag"
	"skybert.net/ytop/pkg"
)

// Populated at build time
var Version = "dev"

type (
	refreshMsg []pkg.Process
)

const (
	headerHeight          = 5
	updateIntervalSeconds = 2
)

var bgColour string
var bgHeaderColour string
var bgSelColour string
var fgColour string
var fgHeaderColour string
var fgSelColour string
var humanSizes bool
var showVersion bool
var simpleView bool

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
	pflag.StringVar(&fgHeaderColour, "header-fg", "#000000", "Header foreground colour")
	pflag.StringVar(&bgHeaderColour, "header-bg", "#06c993", "Header background colour")
	pflag.StringVar(&bgColour, "bg", "#222235", "Background colour")
	pflag.StringVar(&fgColour, "fg", "#b8c0d4", "Foreground colour")
}

type model struct {
	conf        pkg.YTopConf
	height      int
	humanSizes  bool
	processes   []pkg.Process
	searchShow  bool
	searchQuery string
	searchInput textinput.Model
	simpleView  bool
	sortKey     pkg.SortKey
	table       *table.Table
	width       int
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

func main() {
	f, err := tea.LogToFile("debug.log", "ytop")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Couldn't close file: %v\n", err)
		}
	}()

	pflag.Parse()
	if showVersion {
		fmt.Printf("ytop version: %v\n", Version)
		os.Exit(0)
	}

	log.Println("foo")

	t := table.New().Headers(tableHeaders(simpleView)...)
	m := model{
		conf: pkg.YTopConf{
			Foreground:         fgColour,
			Background:         fgColour,
			HeaderForeground:   fgHeaderColour,
			HeaderBackground:   bgHeaderColour,
			SelectedForeground: fgSelColour,
			SelectedBackground: bgSelColour,
			SimpleView:         simpleView,
		},
		table:       t,
		humanSizes:  humanSizes,
		simpleView:  simpleView,
		searchInput: searchInput(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("%v: %v\n", "There was an error", err)
		fmt.Printf("%v\n", debug.Stack())
		os.Exit(1)
	}
}
