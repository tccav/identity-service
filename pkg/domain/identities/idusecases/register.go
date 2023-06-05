package idusecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type RegisterUseCase struct {
	repository    identities.StudentsRepository
	eventProducer identities.StudentsProducer
}

func NewRegisterUseCase(repository identities.StudentsRepository, eventProducer identities.StudentsProducer) RegisterUseCase {
	return RegisterUseCase{
		repository:    repository,
		eventProducer: eventProducer,
	}
}

func (r RegisterUseCase) RegisterStudent(ctx context.Context, input identities.RegisterStudentInput) (string, error) {
	_, err := uuid.Parse(input.CourseID)
	if err != nil {
		return "", fmt.Errorf("%w: %s", identities.ErrInvalidCourseID, err)
	}

	student, err := entities.NewStudent(input.ID, input.Secret, input.Name, input.CPF, input.Email, input.BirthDate)
	if err != nil {
		return "", err
	}

	err = r.repository.CreateStudent(ctx, student)
	if err != nil {
		return "", err
	}

	// FIXIT: if this operations fails, the previous will not be reverted. Think in a transaction like approach to solve it
	err = r.eventProducer.ProduceStudentRegistered(ctx, student, input.CourseID)
	if err != nil {
		return "", err
	}

	return student.ID, nil
}
