package idusecases

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
	"github.com/tccav/identity-service/pkg/gateways/kafka"
	"github.com/tccav/identity-service/pkg/gateways/kafka/kfixtures"
	"github.com/tccav/identity-service/pkg/gateways/postgres"
	"github.com/tccav/identity-service/pkg/gateways/postgres/pgfixtures"
)

func TestRegisterUseCase_RegisterStudent(t *testing.T) {
	t.Parallel()

	validInput := identities.RegisterStudentInput{
		ID:        "201320509911",
		Name:      "Pedro Lopes",
		Secret:    "secret_password",
		CPF:       "11111111030",
		Email:     "plopes@ol.com",
		BirthDate: "1994-03-19",
		CourseID:  uuid.NewString(),
	}

	tt := []struct {
		name    string
		input   identities.RegisterStudentInput
		want    string
		wantErr error
	}{
		{
			name:  "should register student",
			input: validInput,
			want:  validInput.ID,
		},
		{
			name:    "should fail due to invalid country id",
			input:   identities.RegisterStudentInput{CourseID: "invalid"},
			wantErr: identities.ErrInvalidCourseID,
		},
		{
			name: "should fail due to invalid cpf",
			input: identities.RegisterStudentInput{
				ID:       "201126811599",
				CourseID: uuid.NewString(),
				CPF:      "a12345",
			},
			wantErr: entities.ErrInvalidCPF,
		},
		{
			name: "should fail due to invalid email",
			input: identities.RegisterStudentInput{
				ID:       "201126811599",
				CourseID: uuid.NewString(),
				CPF:      "11111111030",
				Email:    "j.com",
			},
			wantErr: entities.ErrInvalidEmail,
		},
		{
			name: "should fail due to invalid email",
			input: identities.RegisterStudentInput{
				ID:        "201126811599",
				CourseID:  uuid.NewString(),
				CPF:       "11111111030",
				Email:     "j@ol.com",
				BirthDate: "1234",
			},
			wantErr: entities.ErrInvalidBirthDate,
		},
	}
	for _, testCase := range tt {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// prepare
			ctx := context.Background()

			dbConn := pgfixtures.NewDB(t)
			repository := postgres.NewStudentsRepository(dbConn)

			kClient := kfixtures.NewKafkaClient(t)
			producer := kafka.NewProducer(kClient)
			eventsProducer := kafka.NewStudentsProducer(producer)

			r := NewRegisterUseCase(repository, eventsProducer)

			// test
			got, err := r.RegisterStudent(ctx, tc.input)

			// assert
			if tc.wantErr != nil {
				assert.Empty(t, got)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NotEmpty(t, got)
				assert.NoError(t, err)
			}
		})
	}
}
