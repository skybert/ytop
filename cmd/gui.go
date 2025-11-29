package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"skybert.net/ytop/pkg"
)

func ViewHeader(sortKey pkg.SortKey, humanSizes bool) string {
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
			sortKey.String(),
		))
	cpu := helpStyle.Render("CPU: ") + CPUSummary()
	memory := helpStyle.Render("Memory: ") + MemorySummary(humanSizes)
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

func (m *model) createTable() table.Model {
	columns := TableColumns(m.simpleView, m.width)
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(false).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(m.conf.SelectedForeground)).
		Background(lipgloss.Color(m.conf.SelectedBackground)).
		Bold(false)
	t.SetStyles(s)

	return t
}
