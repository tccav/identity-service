package identities

import (
	"context"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

type StudentsRepository interface {
	CreateStudent(ctx context.Context, student entities.Student) error
}
