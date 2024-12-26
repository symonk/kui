package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/symonk/kui/internal/kafka"
)

// Frontend encapsulates multiple composite views.
// and is responsible for routing and delegating
// updates to nested models.
// Frontend tracks a stack of models where model[-1]
// is the current model 'in view'.
type Frontend struct {
	client     *kafka.Client
	connected  bool
	modelStack []tea.Model
}

func New(client *kafka.Client) *Frontend {
	return &Frontend{
		client:     client,
		connected:  false,
		modelStack: []tea.Model{&Menu{}},
	}

}

// Init establishes synchroous connectivity to the kafka brokers.
func (f *Frontend) Init() tea.Cmd {
	return kafkaConnectionCommand(f.client)
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
			f.PushModel(&Menu{})
		}
	}
	return f.modelStack[len(f.modelStack)-1].Update(msg)
}

// PopModel closes the view of a particular model, returning the
// view to the one opened previously.
func (f *Frontend) PopModel() {
	if len(f.modelStack) > 1 {
		f.modelStack = f.modelStack[:len(f.modelStack)-1]
	}
}

// PushModel puts a new model into the view.
func (f *Frontend) PushModel(m tea.Model) {
	f.modelStack = append(f.modelStack, m)
}

// View displays the currently active model.
func (f *Frontend) View() string {
	if f.connected == false {
		c := &Connector{}
		f.PushModel(c)
	}
	return f.modelStack[len(f.modelStack)-1].View()
}
