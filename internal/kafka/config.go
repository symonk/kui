package kafka

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// FileToKafkaMap accepts the resolved path to the librdkafka
// file and parses all the values from it in preparation for
// instantiating a kafka admin client (or producer/consumer)
// in future.
//
// Note: Various librdkafka options are only applicable to
// (P) || (C) || (A) and not all.  The underlying client will
// ignore such options that are not sufficient for the type we
// have instantiated.
func FileToKafkaMap(path string) (kafka.ConfigMap, error) {
	c := kafka.ConfigMap{}
	f, err := os.Open(path)
	if err != nil {
		return c, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		k, v, ok := strings.Cut(l, "=")
		if !ok {
			return nil, fmt.Errorf("line in config missing =: %s", l)
		}
		c.SetKey(k, v)
	}
	if s.Err() != nil {
		return c, err
	}

	return c, nil
}
