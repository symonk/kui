package gui

import tea "github.com/charmbracelet/bubbletea"

type DebugView struct {
}

func (d *DebugView) View() string {
	return "DebugView"
}

func (d *DebugView) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return d, nil
}

func (d *DebugView) Init() tea.Cmd {
	return nil
}
