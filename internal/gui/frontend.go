package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/symonk/kui/internal/kafka"
)

type View string

const (
	ConnectingView = "connecting"
	MenuView       = "menu"
)

// Frontend encapsulates multiple composite views.
// and is responsible for routing and delegating
// updates to nested models.
// Frontend tracks a stack of models where model[-1]
// is the current model 'in view'.
type Frontend struct {
	client    *kafka.Client
	connected bool
	router    map[string]tea.Model
	visible   string
}

func New(client *kafka.Client) *Frontend {
	return &Frontend{
		client:    client,
		connected: false,
		router:    map[string]tea.Model{ConnectingView: NewConnector(client)},
		visible:   ConnectingView,
	}

}

// Init establishes synchroous connectivity to the kafka brokers.
func (f *Frontend) Init() tea.Cmd {
	return f.router[f.visible].Init()
}

func (f *Frontend) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return f, tea.Quit
		}
	case connMessage:
		f.connected = bool(msg)
		if !f.connected {
			f.visible = ConnectingView
		}
	}
	return f.router[f.visible].Update(msg)
}

// View displays the currently active model.
func (f *Frontend) View() string {
	if f.connected {
		f.visible = MenuView
	}
	return f.router[f.visible].View()
}
