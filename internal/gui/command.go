package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/symonk/kui/internal/kafka"
)

type connMessage struct {
	ok  bool
	err error
}

// KafkaConnectedCmd establishes connectivity to the kafka
// cluster.  This state is ran on init in order to establish
// the connection.
func kafkaConnectionCommand(client *kafka.Client) tea.Cmd {
	return func() tea.Msg {
		if err := client.WaitForBrokerConnection(); err != nil {
			return connMessage{ok: false, err: err}
		}
		return connMessage{ok: true, err: nil}
	}
}
