package main

// Models used by ytop

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2/table"
	"skybert.net/ytop/pkg"
)

type inputType int

const (
	searchInput inputType = iota
	killInput   inputType = iota
	signalInput inputType = iota
)

type model struct {
	conf       pkg.YTopConf
	height     int
	humanSizes bool
	input      textinput.Model
	inputType  inputType
	inputQuery string
	inputShow  bool
	processes  []pkg.Process
	simpleView bool
	sortKey    pkg.SortKey
	table      *table.Table
	width      int
}
