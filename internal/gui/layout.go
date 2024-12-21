// package gui exposes the layout of the application
package gui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
)

const (
	padding  = 2
	maxWidth = 120
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

// AppModel encapsulates the frontend state of the application.
type AppModel struct {
	client   *kafka.Client
	progress progress.Model
}

// New returns a new pointer to an AppModel instance configured for the
// application.
func New(client *kafka.Client) *AppModel {
	return &AppModel{
		client:   client,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (a *AppModel) Init() tea.Cmd {
	return tickCommand()
}

func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return a, tea.Quit

		}
	case tickMsg:
		if a.progress.Percent() == 1.0 {
			return a, nil
		}
		cmd := a.progress.IncrPercent(0.2)
		return a, tea.Batch(tickCommand(), cmd)
	case progress.FrameMsg:
		mod, cmd := a.progress.Update(msg)
		a.progress = mod.(progress.Model)
		return a, cmd
	default:
		return a, nil
	}
	return a, nil
}

type tickMsg time.Time

func (a *AppModel) View() string {
	pad := strings.Repeat(" ", padding)
	return lipgloss.NewStyle().BorderStyle(lipgloss.OuterHalfBlockBorder()).Render("connecting to brokers..\n\n" +
		pad + a.progress.View() + "\n\n" +
		pad + style.Render("press q to quit\n\n"))
}

func tickCommand() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
