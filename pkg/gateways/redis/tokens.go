package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type TokensRepository struct {
	client *redis.Client
}

func NewTokensRepository(client *redis.Client) TokensRepository {
	return TokensRepository{
		client: client,
	}
}

func (t TokensRepository) Register(ctx context.Context, token entities.Token) error {
	_, err := t.client.SetNX(ctx, parseTokenKey(token.ID), token.Hash, time.Until(token.ExpirationDate)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (t TokensRepository) GetHash(ctx context.Context, id string) (string, error) {
	hash, err := t.client.Get(ctx, parseTokenKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", identities.ErrTokenNotEmitted
		}
		return "", err
	}
	return hash, nil
}

func parseTokenKey(tokenID string) string {
	const tokenKeyTpl = "token:%s"

	return fmt.Sprintf(tokenKeyTpl, tokenID)
}
