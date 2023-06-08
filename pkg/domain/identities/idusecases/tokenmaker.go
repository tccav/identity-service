package idusecases

import (
	"context"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

const (
	defaultTokenDuration = 3 * time.Hour
	defaultIssuer        = "aol"
)

type jwtTokenMaker struct {
	secret     string
	repository identities.TokenRegistererRepository
}

func (m jwtTokenMaker) createToken(ctx context.Context, input createTokenInput) (entities.Token, error) {
	now := time.Now().UTC()
	token := entities.NewToken(input.userID, now.Add(defaultTokenDuration))

	t, err := jwt.NewBuilder().
		JwtID(token.ID).
		Expiration(token.ExpirationDate).
		IssuedAt(now).
		Issuer(defaultIssuer).
		Subject(input.userID).Build()
	if err != nil {
		return entities.Token{}, err
	}

	hash, err := jwt.Sign(t, jwt.WithKey(jwa.HS256, []byte(m.secret)))
	if err != nil {
		return entities.Token{}, err
	}

	token.Hash = string(hash)

	err = m.repository.RegisterToken(ctx, token)
	if err != nil {
		return entities.Token{}, err
	}

	return token, nil
}
