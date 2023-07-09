package kfixtures

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var (
	kafkaURL = "localhost:9094"
)

func NewKafkaClient(t *testing.T) *kgo.Client {
	t.Helper()

	ctx := context.Background()

	client, err := kgo.NewClient(kgo.SeedBrokers(kafkaURL))
	require.NoError(t, err)

	retryCount := 5
	timer := time.NewTicker(5 * time.Second)
	for range timer.C {
		retryCount--
		err = client.Ping(ctx)
		if retryCount == 0 || err == nil {
			timer.Stop()
			break
		}
	}

	require.NoError(t, err)

	admClient := kadm.NewClient(client)

	_, err = admClient.CreateTopics(ctx, 1, 1, nil, "identity.cdc.students.0")
	require.NoError(t, err)

	return client
}
