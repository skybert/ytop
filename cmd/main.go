package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
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
