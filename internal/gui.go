package internal

import "github.com/charmbracelet/bubbles/table"

func TableColumns(
	simpleView bool,
	nameWidth int,
	cmdWidth int) []table.Column {
	columns := []table.Column{
		{Title: "PID", Width: 7},
		{Title: "RSS", Width: 10},
		{Title: "%CPU", Width: 5},
	}

	if simpleView {
		columns = append(columns, table.Column{Title: "NAME", Width: nameWidth + cmdWidth})
	} else {
		columns = append(columns, table.Column{Title: "NAME", Width: nameWidth})
		columns = append(columns, table.Column{Title: "COMMAND", Width: cmdWidth})
	}

	return columns
}
