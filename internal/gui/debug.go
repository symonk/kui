package gui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var debugTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

var cols = []table.Column{
	{Title: "Level", Width: 10},
	{Title: "Message", Width: 70},
}

// dummyRows is for testing purposes for now, eventually we can pipe the
// io.writer and io.reader to supply a stream of debug logs visibile on
// the frontend.
var dummyRows = []table.Row{
	{"INFO", "The quick brown fox jumps over the lazy dog"},
	{"DEBUG", "The yellow cat jumps over the lazy dog"},
	{"ERROR", "The white rabbit jumps over the lazy dog"},
	{"WARN", "The buffalo went to the moon"},
	{"INFO", "The quick brown fox jumps over the lazy dog"},
	{"DEBUG", "The yellow cat jumps over the lazy dog"},
	{"ERROR", "The white rabbit jumps over the lazy dog"},
	{"WARN", "The buffalo went to the moon"},
	{"INFO", "The quick brown fox jumps over the lazy dog"},
	{"DEBUG", "The yellow cat jumps over the lazy dog"},
	{"ERROR", "The white rabbit jumps over the lazy dog"},
	{"WARN", "The buffalo went to the moon"},
}

// DebugView is a live trace of the current logging
// The vision for DebugView is to encapsulate all logging
// into two sections (kui and core kafka debug) and display
// them in a table.  These are also written to a *os.File
// somewhere on disk and potentially tailed back in here
// or perhaps using some multi writer with IO.Pipe to connect
// such writer<->readers together.
//
// Smart selection of rows, which includes scrolling capabilities
// aswell as the ability to select and cause a prompt of the
// entire message in pretty printed json.
//
// No such plans to make the actual json inspectable, a full
// view will suffice for now.
//
// Lastly, additional styling should be added for consistency throughout
// all tabs.
type DebugView struct {
	table table.Model
}

func NewDebugView() *DebugView {
	dv := &DebugView{
		table: table.New(
			table.WithFocused(true),
			table.WithHeight(10),
			table.WithWidth(80),
			table.WithColumns(cols),
			table.WithRows(dummyRows),
		),
	}
	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	dv.table.SetStyles(s)
	return dv
}

func (d *DebugView) View() string {
	return debugTableStyle.Render(d.table.View()) + "\n" + d.table.HelpView() + "\n"
}

func (d *DebugView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if d.table.Focused() {
				d.table.Blur()
			} else {
				d.table.Focus()
			}
		case "q":
			return d, tea.Quit
		case "enter":
		}
	}
	mod, cmd := d.table.Update(msg)
	d.table = mod
	return d, cmd
}

func (d *DebugView) Init() tea.Cmd {
	return nil
}
