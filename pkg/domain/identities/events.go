package identities

import (
	"context"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

type StudentsProducer interface {
	ProduceStudentRegistered(ctx context.Context, student entities.Student, courseID string) error
}
