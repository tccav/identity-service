package kafka

import (
	"context"

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
		topic: "foo",
		event: event{
			ID:   uuid.NewString(),
			Type: "student_registered",
			Payload: studentRegisteredPayload{
				StudentID: student.ID,
				Name:      student.Name,
				CPF:       student.CPF,
				Email:     student.Email,
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
