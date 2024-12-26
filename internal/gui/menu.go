package gui

import tea "github.com/charmbracelet/bubbletea"

type Menu struct{}

func (m *Menu) View() string {
	return "connected"
}

func (m *Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Menu) Init() tea.Cmd {
	return nil
}
