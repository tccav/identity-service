package kafka

type event struct {
	ID      string `json:"event_id"`
	Type    string `json:"event_type"`
	Payload any    `json:"payload"`
}

type studentRegisteredPayload struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	CourseID  string `json:"course_id"`
}
