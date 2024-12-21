// package kafka exposes a kafka admin client
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client struct {
	client *kafka.AdminClient
}

// New instantiates a new instance of the kafka client and connects
// to the brokers specified in the bootstrap servers.
func New() (*Client, error) {
	c, err := kafka.NewAdminClient(&kafka.ConfigMap{})
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}
