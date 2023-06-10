package identities

import (
	"context"
	"errors"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

//go:generate moq -out idmocks/mock_usecases.go -pkg idmocks . RegisterUseCases AuthenticationUseCases

var (
	ErrInvalidCourseID      = errors.New("invalid course id")
	ErrStudentAlreadyExists = errors.New("student already exists")
	ErrStudentNotFound      = errors.New("student not found")

	ErrEmptyStudentID   = errors.New("empty student id was sent")
	ErrEmptySecret      = errors.New("empty secret was sent")
	ErrSecretsDontMatch = errors.New("secrets don't match")
	ErrEmptyToken       = errors.New("token hash informed is empty")
	ErrTokenExpired     = errors.New("token expired")
	ErrMalformedToken   = errors.New("malformed token")
	ErrTokenNotEmitted  = errors.New("informed token was not emitted by this app")
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

type AuthenticateStudentInput struct {
	StudentID     string
	StudentSecret string
}

type AuthenticationUseCases interface {
	AuthenticateStudent(ctx context.Context, input AuthenticateStudentInput) (entities.Token, error)
	VerifyAuth(ctx context.Context, hash string) error
}
