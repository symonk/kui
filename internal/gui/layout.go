// package gui exposes the layout of the application
package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/symonk/kui/internal/kafka"
)

// AppModel encapsulates the frontend state of the application.
type AppModel struct {
	client   *kafka.Client
	choices  []string
	cursor   int
	selected map[int]struct{}
}

// New returns a new pointer to an AppModel instance configured for the
// application.
func New(client *kafka.Client) *AppModel {
	choices := make([]string, 0)
	topics := client.FetchTopics()
	for _, topic := range topics {
		choices = append(choices, topic.Topic)
	}
	return &AppModel{
		client:   client,
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

func (a *AppModel) Init() tea.Cmd {
	return nil
}

func (a *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return a, tea.Quit

		}

	}
	return a, nil
}

func (a *AppModel) View() string {
	in := "What would you like to do?"
	for i, choice := range a.choices {
		cursor := " "
		if a.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := a.selected[i]; ok {
			checked = "x"
		}
		in += "\n" + cursor + " [" + checked + "] " + choice
	}
	in += "\nPress q to quit.\n"
	return in
}
