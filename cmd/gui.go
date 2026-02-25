package main

import (
	"fmt"
	"log"
	"time"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
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

func (m *model) createTable() *table.Table {
	columns := tableHeaders(m.simpleView)
	log.Printf("hdr bg %v\n", m.conf.HeaderBackground)
	log.Printf("hdr fg %v\n", m.conf.HeaderForeground)

	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.
		Background(lipgloss.Color(m.conf.HeaderBackground)).
		Foreground(lipgloss.Color(m.conf.HeaderForeground)).
		Bold(true)

	t := table.New().
		Headers(columns...).
		BorderTop(false).
		BorderColumn(false).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
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
