package idusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type RegisterUseCase struct {
	repository    identities.StudentsRegistererRepository
	eventProducer identities.StudentsProducer
	tracer        trace.Tracer
}

func NewRegisterUseCase(repository identities.StudentsRegistererRepository, eventProducer identities.StudentsProducer) RegisterUseCase {
	return RegisterUseCase{
		repository:    repository,
		eventProducer: eventProducer,
		tracer:        otel.Tracer(tracerName),
	}
}

func (r RegisterUseCase) RegisterStudent(ctx context.Context, input identities.RegisterStudentInput) (string, error) {
	ctx, span := r.tracer.Start(ctx, "RegisterUseCase.RegisterStudent")
	defer span.End()

	_, err := uuid.Parse(input.CourseID)
	if err != nil {
		err = fmt.Errorf("%w: %s", identities.ErrInvalidCourseID, err)
		span.RecordError(err)
		return "", err
	}

	student, err := entities.NewStudent(input.ID, input.Secret, input.Name, input.CPF, input.Email, input.BirthDate)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	err = r.repository.CreateStudent(ctx, student)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	// FIXIT: if this operations fails, the previous will not be reverted. Think in a transaction like approach to solve it
	err = r.eventProducer.ProduceStudentRegistered(ctx, student, input.CourseID)
	if err != nil {
		return "", err
	}

	return student.ID, nil
}
