// package kafka exposes a kafka admin client
package kafka

import (
	_ "github.com/charmbracelet/bubbletea"
	_ "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client struct {
}

func New() *Client {
	return &Client{}
}
