package idusecases

import (
	"context"
	"testing"
	"time"

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

		s := NewStudentJWTAuthenticator(studentsRepository, tokensRepository, "test")

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
			name: "should fail due to empty secret",
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

			s := NewStudentJWTAuthenticator(studentsRepository, nil, "test")

			got, err := s.AuthenticateStudent(ctx, tc.input)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Empty(t, got)
		})
	}
}
