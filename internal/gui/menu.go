package gui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
	"golang.org/x/term"
)

var menuStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

var activeTabBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      " ",
	Left:        "|",
	Right:       "|",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┘",
	BottomRight: "└",
}

var tabBorder = lipgloss.Border{
	Top:         "─",
	Bottom:      "─",
	Left:        "│",
	Right:       "│",
	TopLeft:     "╭",
	TopRight:    "╮",
	BottomLeft:  "┴",
	BottomRight: "┴",
}

var highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
var tabStyle = lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(highlight).Padding(0, 1)
var activeTab = tabStyle.Border(activeTabBorder, true)
var tabGap = tabStyle.BorderTop(false).BorderLeft(false).BorderRight(false)

// Menu is the core view of kui.
type Menu struct {
	client    *kafka.Client
	status    *StatusBar
	connected bool
	width     int
	height    int
}

func NewMenu(client *kafka.Client) *Menu {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return &Menu{client: client, status: NewStatusBar(client), width: w, height: h}
}

func (m *Menu) View() string {
	var doc strings.Builder
	row := lipgloss.JoinHorizontal(lipgloss.Top, activeTab.Render("Overview"), tabStyle.Render("Topics"), tabStyle.Render("Consumer Groups"), tabStyle.Render("Settings"))
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(row))))
	row = lipgloss.JoinHorizontal(lipgloss.Top, row, gap)
	_, _ = doc.WriteString(row + "\n\n")
	return doc.String() + m.status.View()
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case connMessage:
		m.status.Update(msg)
	}
	return m, nil
}

func (m *Menu) Init() tea.Cmd {
	return nil
}
