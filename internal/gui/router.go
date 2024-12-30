package gui

import (
	"log/slog"

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
	logger  *slog.Logger
}

func New(client *kafka.Client, logger *slog.Logger, logFile string) *Router {
	return &Router{
		client:  client,
		router:  map[string]tea.Model{ConnectingView: NewConnector(client, logFile)},
		visible: ConnectingView,
		logger:  logger,
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
