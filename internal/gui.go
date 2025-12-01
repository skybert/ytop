package internal

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func TableColumns(simpleView bool, totalWidth int) []table.Column {
	pidWidth := 7
	rssWidth := 10
	cpuWidth := 5
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "RSS", Width: 10},
		{Title: "%CPU", Width: 5},
	}

	if simpleView {
		nameWidth := totalWidth - pidWidth - rssWidth - cpuWidth
		columns = append(columns, table.Column{Title: "NAME", Width: nameWidth})
	} else {
		nameWidth := 10
		cmdWidth := totalWidth - pidWidth - rssWidth - cpuWidth - nameWidth
		columns = append(columns, table.Column{Title: "NAME", Width: nameWidth})
		columns = append(columns, table.Column{Title: "COMMAND", Width: cmdWidth})
	}

	return columns
}

func CreateTable(simpleView bool, totalWidth int) table.Model {
	columns := TableColumns(simpleView, totalWidth)
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
