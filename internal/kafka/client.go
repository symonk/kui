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
func New(cfgMap kafka.ConfigMap) (*Client, error) {
	c, err := kafka.NewAdminClient(&cfgMap)
	if err != nil {
		return nil, err
	}
	client := &Client{client: c}
	client.WaitForBrokerConnection()
	return client, nil
}

func (c *Client) WaitForBrokerConnection() error {
	if _, err := c.client.GetMetadata(nil, false, 5000); err != nil {
		return err
	}
	return nil
}
