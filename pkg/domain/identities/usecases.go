package identities

import (
	"context"
	"errors"
)

//go:generate moq -out mocks/mock_usecases.go -pkg idmocks . RegisterUseCases

var (
	ErrInvalidCourseID      = errors.New("invalid course id")
	ErrStudentAlreadyExists = errors.New("student already exists")
)

type RegisterStudentInput struct {
	ID        string
	Name      string
	Secret    string
	CPF       string
	Email     string
	BirthDate string
	CourseID  string
}

type RegisterUseCases interface {
	RegisterStudent(ctx context.Context, input RegisterStudentInput) (string, error)
}
