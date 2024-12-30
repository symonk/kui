package gui

import (
	tea "github.com/charmbracelet/bubbletea"
	confluentKafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
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
			return connMessage{ok: true, err: err}
		}
		return connMessage{ok: true, err: nil}
	}
}

type MetaData struct {
	meta *confluentKafka.Metadata
	err  error
}

// kafkaMetaDataCommand is asynchronously dispatched by bubble tea and performs
// a cluster meta data fetch to kafka.  The meta data returned includes information
// on the cluster, aswell as (all) topics, inclusive of internal ones.
func kafkaMetaDataCommand(client *kafka.Client) tea.Cmd {
	return func() tea.Msg {
		var m MetaData
		meta, err := client.FetchMetaData()
		m.meta = meta
		m.err = err
		return m
	}
}
