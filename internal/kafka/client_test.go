package kafka

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

func TestContainersExample(t *testing.T) {
	t.Skip("work in progress")
	ctx := context.Background()
	container, err := kafka.Run(ctx, "confluentinc/confluent-local:7.5.0", kafka.WithClusterID("test-cluster"))
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 3)
	assert.True(t, true)
}
