package rfixtures

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func NewDB(t *testing.T) *redis.Client {
	t.Helper()

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	_, err := client.Ping(context.Background()).Result()
	require.NoError(t, err)

	return client
}
