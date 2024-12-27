package gui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
)

var statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

var statusColourMap = map[string]lipgloss.Color{
	"connected":    lipgloss.Color("#00FF00"),
	"disconnected": lipgloss.Color("#FF0000"),
}

// StatusBar is a simple status bar that displays connectivity
// information.
type StatusBar struct {
	client *kafka.Client
	mode   string
}

func NewStatusBar(client *kafka.Client) *StatusBar {
	return &StatusBar{client: client, mode: "disconnected"}
}

func (s *StatusBar) Init() tea.Cmd {
	return kafkaConnectionCommand(s.client)
}

func (s *StatusBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connMessage:
		if msg.ok {
			s.mode = "connected"
			return s, nil
		}
		s.mode = "disconnected"
		return s, nil
	}
	return s, nil
}

func (s *StatusBar) View() string {
	return lipgloss.NewStyle().Foreground(statusColourMap[s.mode]).Render(fmt.Sprintf("STATUS: %s", s.mode))
}
