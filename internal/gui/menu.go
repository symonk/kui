package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/symonk/kui/internal/kafka"
)

// Menu is the core view of kui.
type Menu struct {
	client    *kafka.Client
	connected bool
}

func NewMenu(client *kafka.Client) *Menu {
	return &Menu{client: client}
}

func (m *Menu) View() string {
	return "connected, press q to quit"
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Menu) Init() tea.Cmd {
	return nil
}
