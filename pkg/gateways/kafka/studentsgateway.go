package kafka

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/tccav/identity-service/pkg/domain/entities"
)

type StudentsGateway struct {
	producer Producer
}

func NewStudentsProducer(producer Producer) StudentsGateway {
	return StudentsGateway{
		producer: producer,
	}
}

func (g StudentsGateway) ProduceStudentRegistered(ctx context.Context, student entities.Student, courseID string) error {
	err := g.producer.produce(ctx, produceInput{
		topic: "identity.cdc.students.0",
		event: event{
			ID:   uuid.NewString(),
			Type: "student_registered",
			Payload: studentRegisteredPayload{
				StudentID: student.ID,
				Name:      student.Name,
				CPF:       student.CPF,
				Email:     student.Email,
				BirthDate: student.BirthDate.Format(time.DateOnly),
				CourseID:  courseID,
			},
		},
	},
	)
	if err != nil {
		return err
	}

	return nil
}
