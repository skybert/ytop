package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2/table"
)

type model struct {
	table *table.Table
}

// Init implements [tea.Model].
func (m model) Init() tea.Cmd {
	return nil
}

// Update implements [tea.Model].
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table = m.table.Width(msg.Width)
		m.table = m.table.Height(msg.Height)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, cmd
}

// View implements [tea.Model].
func (m model) View() tea.View {
	v := tea.NewView("\n" + m.table.String() + "\n")
	v.AltScreen = true
	return v
}

func main() {
	headers := []string{"Name", "Weapon", "Song"}
	rows := [][]string{
		{"Aragorn", "Sword", "Et Eärello Endorenna utúlien Sinome maruvan ar Hildinyar tenn' Ambar-metta!"},
		{"Arwen", "Bow", "With a sigh You turn away With a deepening heart No words to say"},
	}
	t := table.New().
		Headers(headers...).
		Rows(rows...)
	if _, err := tea.NewProgram(model{t}).Run(); err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
}
