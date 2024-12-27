// package kafka exposes a kafka admin client
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client struct {
	client  *kafka.AdminClient
	timeout int
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

// FetchTopicsMetaData fetches information on all topics.
// including their partition setups.
func (c *Client) FetchTopics() []kafka.TopicMetadata {
	meta, err := c.client.GetMetadata(nil, false, 2000)
	if err != nil {
		return nil
	}
	t := make([]kafka.TopicMetadata, len(meta.Topics))
	for _, v := range meta.Topics {
		t = append(t, v)
	}
	return t
}

func (c *Client) FetchMetaData() (*kafka.Metadata, error) {
	return c.client.GetMetadata(nil, false, 2000)
}

func (c *Client) WaitForBrokerConnection() error {
	if _, err := c.client.GetMetadata(nil, false, 1000); err != nil {
		return err
	}
	return nil
}
