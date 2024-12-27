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

// Router encapsulates multiple composite views.
// and is responsible for routing and delegating
// updates to nested models.
// Router tracks a stack of models where model[-1]
// is the current model 'in view'.
type Router struct {
	client  *kafka.Client
	router  map[string]tea.Model
	visible string
}

func New(client *kafka.Client) *Router {
	return &Router{
		client:  client,
		router:  map[string]tea.Model{ConnectingView: NewConnector(client)},
		visible: ConnectingView,
	}

}

// Init establishes synchroous connectivity to the kafka brokers.
func (f *Router) Init() tea.Cmd {
	return f.router[f.visible].Init()
}

func (f *Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f.router[f.visible].Update(msg)
}

// View displays the currently active model.
func (f *Router) View() string {
	return f.router[f.visible].View()
}