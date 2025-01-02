package gui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/symonk/kui/internal/kafka"
	"github.com/symonk/kui/internal/terminal"
)

var (

	// connectorStyle is the styling for the connector border.
	// this is placed in the center of the users view.
	connectorStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)

	// connectingMessageStyle is the styling for the current message
	// including any kafka errors etc.
	connectingMessageStyle = lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Center)

	// logFileMessageStyle is the styling for the log file text
	logFileMessageStyle = lipgloss.NewStyle().
				Align(lipgloss.Center)

	// quitMessageStyle is the styling for the quit message
	quitMessageStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				Foreground(lipgloss.Color("#626262"))
)

type Connector struct {
	client    *kafka.Client
	progress  progress.Model
	message   string
	brokers   []string
	connected bool
	logFile   string
	width     int
	height    int
}

// NewConnector returns an instance of the Connector model.  The
// Connector model instance is displayed on initial connection or
// if the connection is dropped while the program is running.
func NewConnector(client *kafka.Client, logFile string) *Connector {
	width, height := terminal.Size()
	return &Connector{
		client:   client,
		progress: progress.New(progress.WithDefaultGradient(), progress.WithWidth(width/2)),
		message:  "connecting to brokers",
		brokers:  make([]string, 0),
		logFile:  logFile,
		width:    width,
		height:   height,
	}
}

// View returns the string responsible for drawing the view of the
// connector window and any of it's internal composite components
// concatennated appropriately.
func (c *Connector) View() string {
	message := connectingMessageStyle.Render(fmt.Sprintf("%s\n\n", c.message))
	filePath := logFileMessageStyle.Render(fmt.Sprintf("log: %s\n", c.logFile))
	quit := quitMessageStyle.Render("press 'ctl+c' to quit")
	progress := c.progress.View() + "\n\n"
	ui := lipgloss.JoinVertical(lipgloss.Center, message, progress, filePath, quit)
	box := lipgloss.Place(c.width, c.height, lipgloss.Center, lipgloss.Center, connectorStyle.Render(ui))
	return box
}

type tickMsg time.Time

func (c *Connector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.updateSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
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
			return NewMenu(c.client), nil
		}
		if !c.connected && c.progress.Percent() >= 0.8 {
			cmd := c.progress.DecrPercent(0.1)
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

// updateSize is responsible for scaling all components
// when sizing events occur.
func (c *Connector) updateSize(scale tea.WindowSizeMsg) {
	c.width, c.height = scale.Width, scale.Height
	c.progress.Width = scale.Width / 4
}

func (c *Connector) Init() tea.Cmd {
	return tea.Batch(kafkaConnectionCommand(c.client), tickCommand())
}

func tickCommand() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
