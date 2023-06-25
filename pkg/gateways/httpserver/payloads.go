package httpserver

import (
	"encoding/json"
	"net/http"
)

type HTTPError struct {
	Code    string `json:"err_code"`
	Message string `json:"message"`
}

var (
	invalidJSON = HTTPError{
		Code:    "identity_service.error.invalid_json",
		Message: "invalid JSON was sent",
	}
	unexpectedError = HTTPError{
		Code:    "identity_service.error.unexpected",
		Message: "Unexpected Error",
	}

	invalidStudentID = HTTPError{
		Code:    "identity_service.error.invalid_student_id",
		Message: "Invalid Student ID was sent",
	}
	invalidCPF = HTTPError{
		Code:    "identity_service.error.invalid_cpf",
		Message: "Invalid CPF was sent",
	}
	invalidEmail = HTTPError{
		Code:    "identity_service.error.invalid_email",
		Message: "Invalid email was sent",
	}
	invalidBirthDate = HTTPError{
		Code:    "identity_service.error.invalid_birth_date",
		Message: "Invalid birth date was sent",
	}
	invalidCourseID = HTTPError{
		Code:    "identity_service.error.invalid_course_id",
		Message: "Invalid course id format was sent",
	}
)

func sendJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		return err
	}
	return nil
}
