package kfixtures

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

const kafkaURL = "localhost:9094"

func NewKafkaClient(t *testing.T) *kgo.Client {
	t.Helper()

	ctx := context.Background()

	client, err := kgo.NewClient(kgo.SeedBrokers(kafkaURL))
	require.NoError(t, err)

	err = client.Ping(ctx)
	require.NoError(t, err)

	admClient := kadm.NewClient(client)

	_, err = admClient.CreateTopics(ctx, 1, 1, nil, "foo")
	require.NoError(t, err)

	return client
}
