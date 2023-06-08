package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

type TokensRepository struct {
	client *redis.Client
}

func NewTokensRepository(client *redis.Client) TokensRepository {
	return TokensRepository{
		client: client,
	}
}

func (t TokensRepository) RegisterToken(ctx context.Context, token entities.Token) error {
	_, err := t.client.SetNX(ctx, parseTokenKey(token.ID), token.Hash, time.Until(token.ExpirationDate)).Result()
	if err != nil {
		return err
	}
	return nil
}

func parseTokenKey(tokenID string) string {
	const tokenKeyTpl = "token:%s"

	return fmt.Sprintf(tokenKeyTpl, tokenID)
}
