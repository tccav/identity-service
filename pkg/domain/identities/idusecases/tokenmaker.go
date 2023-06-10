package idusecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type jwtTokenMaker struct {
	secret     string
	issuer     string
	duration   time.Duration
	repository identities.TokenRegistererRepository
}

func (m jwtTokenMaker) createToken(ctx context.Context, userID string) (entities.Token, error) {
	token, err := m.buildSignedJWT(userID)
	if err != nil {
		return entities.Token{}, err
	}

	err = m.repository.Register(ctx, token)
	if err != nil {
		return entities.Token{}, err
	}

	return token, nil
}

func (m jwtTokenMaker) buildSignedJWT(userID string) (entities.Token, error) {
	now := time.Now().UTC()
	token := entities.NewToken(userID, now.Add(m.duration))

	t, err := jwt.NewBuilder().
		JwtID(token.ID).
		Expiration(token.ExpirationDate).
		IssuedAt(now).
		Issuer(m.issuer).
		Subject(token.UserID).Build()
	if err != nil {
		return entities.Token{}, err
	}

	hash, err := jwt.Sign(t, jwt.WithKey(jwa.HS256, []byte(m.secret)))
	if err != nil {
		return entities.Token{}, err
	}

	token.Hash = string(hash)
	return token, nil
}

func (m jwtTokenMaker) verifyToken(ctx context.Context, hash string) error {
	token, err := jwt.ParseString(
		hash,
		jwt.WithIssuer(m.issuer),
		jwt.WithKey(jwa.HS256, []byte(m.secret)),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired()) {
			return identities.ErrTokenExpired
		}
		return fmt.Errorf("%w: %s", identities.ErrMalformedToken, err)
	}

	storedHash, err := m.repository.GetHash(ctx, token.JwtID())
	if err != nil {
		return err
	}

	if storedHash != hash {
		return identities.ErrTokenNotEmitted
	}

	return nil
}
