package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))

type Connector struct {
	client    *kafka.Client
	progress  progress.Model
	message   string
	brokers   []string
	connected bool
}

func NewConnector(client *kafka.Client) *Connector {
	return &Connector{
		client:   client,
		progress: progress.New(progress.WithDefaultGradient()),
		message:  "connecting to brokers",
		brokers:  make([]string, 0),
	}
}

func (c *Connector) View() string {
	pad := strings.Repeat(" ", 2)
	return "\n" + lipgloss.NewStyle().Bold(true).Padding(1).Width(50).Align(0, 0).
		BorderStyle(lipgloss.RoundedBorder()).Align(lipgloss.Center).Render(fmt.Sprintf("  %s..\n\n", c.message)+
		pad+c.progress.View()+"\n\n"+
		pad+style.Render("Press 'q' to quit")) + "\n"
}

type tickMsg time.Time

func (c *Connector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return c, tea.Quit
		}
	case connMessage:
		if msg.ok {
			c.connected = true
			return c, nil
		}
		c.message = msg.err.Error()
		return c, kafkaConnectionCommand(c.client)
	case tickMsg:
		if c.connected && c.progress.Percent() == 1.0 {
			c.message = "successfully connected"
			c.connected = true
			return NewMenu(c.client), nil
		}
		if !c.connected && c.progress.Percent() >= 0.9 {
			c.message = "connection failing..."
			cmd := c.progress.DecrPercent(0.2)
			return c, tea.Batch(tickCommand(), cmd)
		}
		cmd := c.progress.IncrPercent(0.1)
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
	return tea.Batch(kafkaConnectionCommand(c.client), tickCommand())
}

func tickCommand() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
