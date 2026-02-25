package main

import (
	"os"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	input    textinput.Model
	quitting bool
	query    string
}

func createModel() model {
	ti := textinput.New()
	ti.Placeholder = "Cmd to search for"
	ti.CharLimit = 100
	ti.SetWidth(ti.CharLimit)
	ti.Focus()
	return model{
		input: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.query = m.input.Value()
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.input.VirtualCursor() {
		c = m.input.Cursor()
	}

	s := lipgloss.JoinVertical(
		lipgloss.Top,
		m.input.View(),
		m.footerView(),
	)

	if m.quitting {
		s += "\nGoodbye, ta ta for now!\n"
	}

	v := tea.NewView(s)
	v.Cursor = c
	return v
}

func (m model) footerView() string {
	if m.query != "" {
		return "You searched for: " + m.query
	}

	return ""

}

func main() {
	app := tea.NewProgram(createModel())

	if _, err := app.Run(); err != nil {
		os.Exit(1)
	}
}
