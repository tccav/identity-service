package idusecases

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
	"github.com/tccav/identity-service/pkg/gateways/postgres"
	"github.com/tccav/identity-service/pkg/gateways/postgres/pgfixtures"
	"github.com/tccav/identity-service/pkg/gateways/redis"
	"github.com/tccav/identity-service/pkg/gateways/redis/rfixtures"
)

type config struct {
	secret   string
	issuer   string
	duration time.Duration
}

func (v config) TokenSecret() string {
	return v.secret
}

func (v config) TokenIssuer() string {
	return v.issuer
}

func (v config) TokenDuration() time.Duration {
	return v.duration
}

var validConfig = config{
	secret:   "secret_secret",
	issuer:   "uerj",
	duration: time.Hour,
}

func TestStudentAuthenticator_AuthenticateStudent(t *testing.T) {
	t.Parallel()

	t.Run("should match credentials and create token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		password := "test_password"
		encPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		require.NoError(t, err)

		validStudent := entities.Student{
			ID:        "201116548712",
			Name:      "John Doe",
			Secret:    string(encPassword),
			CPF:       "11111111030",
			Email:     "jdoe@ol.com",
			BirthDate: time.Now().UTC(),
		}

		db := pgfixtures.NewDB(t)
		studentsRepository := postgres.NewStudentsRepository(db)

		err = studentsRepository.CreateStudent(ctx, validStudent)
		require.NoError(t, err)

		rDB := rfixtures.NewDB(t)
		tokensRepository := redis.NewTokensRepository(rDB)

		s := NewStudentJWTAuthenticator(studentsRepository, tokensRepository, validConfig)

		got, err := s.AuthenticateStudent(ctx, identities.AuthenticateStudentInput{
			StudentID:     validStudent.ID,
			StudentSecret: password,
		})

		assert.NoError(t, err)
		assert.Equal(t, validStudent.ID, got.UserID)
		assert.NotEmpty(t, got.ID)
		assert.NotEmpty(t, got.ExpirationDate)
		assert.NotEmpty(t, got.Hash)
	})

	tt := []struct {
		name    string
		input   identities.AuthenticateStudentInput
		wantErr error
	}{
		{
			name: "should fail due to empty student id",
			input: identities.AuthenticateStudentInput{
				StudentSecret: "foo",
			},
			wantErr: identities.ErrEmptyStudentID,
		},
		{
			name: "should fail due to empty student secret",
			input: identities.AuthenticateStudentInput{
				StudentID: "123456789",
			},
			wantErr: identities.ErrEmptySecret,
		},
		{
			name: "should fail because student does not exist",
			input: identities.AuthenticateStudentInput{
				StudentID:     "123456789",
				StudentSecret: "test_password",
			},
			wantErr: identities.ErrStudentNotFound,
		},
	}
	for _, testCase := range tt {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			db := pgfixtures.NewDB(t)
			studentsRepository := postgres.NewStudentsRepository(db)

			s := NewStudentJWTAuthenticator(studentsRepository, nil, validConfig)

			got, err := s.AuthenticateStudent(ctx, tc.input)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Empty(t, got)
		})
	}
}

func TestStudentAuthenticator_VerifyAuth(t *testing.T) {
	t.Parallel()

	t.Run("should verify the auth token with success", func(t *testing.T) {
		t.Parallel()

		// prepare
		ctx := context.Background()

		rDB := rfixtures.NewDB(t)
		tokensRepository := redis.NewTokensRepository(rDB)

		s := NewStudentJWTAuthenticator(nil, tokensRepository, validConfig)

		userID := uuid.NewString()
		token, err := s.createToken(ctx, userID)
		require.NoError(t, err)

		// test
		err = s.VerifyAuth(ctx, token.Hash)

		// assert
		assert.NoError(t, err)
	})

	tt := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "should fail because empty token hash was informed",
			input:   "",
			wantErr: identities.ErrEmptyToken,
		},
		{
			name: "should fail because token issue does not match",
			input: generateToken(t, config{
				issuer:   "ufrj",
				secret:   validConfig.secret,
				duration: validConfig.duration,
			}).Hash,
			wantErr: identities.ErrMalformedToken,
		},
		{
			name: "should fail because token secret does not match",
			input: generateToken(t, config{
				issuer:   validConfig.issuer,
				secret:   "segredo",
				duration: validConfig.duration,
			}).Hash,
			wantErr: identities.ErrMalformedToken,
		},
		{
			name: "should fail because token is expired",
			input: generateToken(t, config{
				issuer:   validConfig.issuer,
				secret:   validConfig.secret,
				duration: -validConfig.duration,
			}).Hash,
			wantErr: identities.ErrTokenExpired,
		},
		{
			name: "should fail because token was not emitted by us",
			input: generateToken(t, config{
				issuer:   validConfig.issuer,
				secret:   validConfig.secret,
				duration: validConfig.duration,
			}).Hash,
			wantErr: identities.ErrTokenNotEmitted,
		},
	}
	for _, testCase := range tt {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			rDB := rfixtures.NewDB(t)
			tokensRepository := redis.NewTokensRepository(rDB)

			s := NewStudentJWTAuthenticator(nil, tokensRepository, validConfig)

			err := s.VerifyAuth(ctx, tc.input)

			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func generateToken(t *testing.T, config config) entities.Token {
	t.Helper()

	token, err := jwtTokenMaker{
		secret:   config.secret,
		issuer:   config.issuer,
		duration: config.duration,
	}.buildSignedJWT(uuid.NewString())
	require.NoError(t, err)

	return token
}
