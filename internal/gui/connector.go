package gui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

type Connector struct {
	client   *kafka.Client
	progress progress.Model
}

func NewConnector(client *kafka.Client) *Connector {
	return &Connector{
		client:   client,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

func (c *Connector) View() string {
	pad := strings.Repeat(" ", 2)
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).Render("connecting to brokers...\n\n") +
		pad + c.progress.View() + "\n\n" +
		pad + style.Render("Press 'q' to quit")
}

func (c *Connector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return c, tea.Quit

		}
	case time.Time:
		if c.progress.Percent() == 1.0 {
			return c, nil
		}
		cmd := c.progress.IncrPercent(0.2)
		return c, tea.Batch(tickCommand(), cmd)
	case progress.FrameMsg:
		mod, cmd := c.progress.Update(msg)
		c.progress = mod.(progress.Model)
		return c, cmd
	default:
		return c, nil
	}
	return c, nil

}

func (c *Connector) Init() tea.Cmd {
	return tickCommand()
}

func tickCommand() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return time.Time(t)
	})
}
