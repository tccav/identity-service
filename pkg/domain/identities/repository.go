package identities

import (
	"context"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

type StudentsRegistererRepository interface {
	CreateStudent(ctx context.Context, student entities.Student) error
}

type StudentListerRepository interface {
	GetStudentSecret(ctx context.Context, id string) (string, error)
}

type TokenRegistererRepository interface {
	Register(ctx context.Context, token entities.Token) error
	GetHash(ctx context.Context, id string) (string, error)
}
