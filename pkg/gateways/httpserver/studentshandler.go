package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type StudentRegisterRequest struct {
	ID        string `json:"id" swaggertype:"string" example:"201210204310"`
	Name      string `json:"name" swaggertype:"string" example:"John Doe"`
	Secret    string `json:"secret" swaggertype:"string" example:"celacanto provoca maremoto"`
	CPF       string `json:"cpf" swaggertype:"string" example:"11111111030"`
	Email     string `json:"email" swaggertype:"string" format:"email" example:"jdoe@ol.com"`
	BirthDate string `json:"birth_date" swaggertype:"string" format:"date" example:"1990-10-18"`
	CourseID  string `json:"course_id" swaggertype:"string" format:"uuidv4" example:"1f6a4d3a-38c7-43fe-9790-2408fe595c93"`
}

type StudentRegisterResponse struct {
	ID string `json:"id" swaggertype:"string" example:"201210204310"`
}

type StudentsHandler struct {
	logger  *zap.Logger
	useCase identities.RegisterUseCases
}

func NewStudentsHandler(useCase identities.RegisterUseCases, logger *zap.Logger) StudentsHandler {
	return StudentsHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterStudent ...
// ShowEntity godoc
// @Summary Register a student
// @Tags Registration
// @Param request body StudentRegisterRequest true "Student creation information"
// @Accept json
// @Produce json
// @Success 201 {object} StudentRegisterResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/identities/students [post]
func (h StudentsHandler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody StudentRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		h.logger.Error("invalid json received", zap.Error(err))
		err = sendJSON(w, http.StatusBadRequest, invalidJSON)
		if err != nil {
			h.logger.Error("failed to send error response", zap.Error(err))
		}
		return
	}

	studentID, err := h.useCase.RegisterStudent(ctx, identities.RegisterStudentInput(reqBody))
	if err != nil {
		h.logger.Error("unable to register user", zap.Error(err))

		var (
			errorPayload HTTPError
			statusCode   int
		)
		switch {
		case errors.Is(err, entities.ErrInvalidStudentID), errors.Is(err, identities.ErrStudentAlreadyExists):
			statusCode = http.StatusBadRequest
			errorPayload = invalidStudentID
		case errors.Is(err, entities.ErrInvalidCPF):
			statusCode = http.StatusBadRequest
			errorPayload = invalidCPF
		case errors.Is(err, entities.ErrInvalidEmail):
			statusCode = http.StatusBadRequest
			errorPayload = invalidEmail
		case errors.Is(err, entities.ErrInvalidBirthDate):
			statusCode = http.StatusBadRequest
			errorPayload = invalidBirthDate
		case errors.Is(err, identities.ErrInvalidCourseID):
			statusCode = http.StatusBadRequest
			errorPayload = invalidCourseID
		default:
			statusCode = http.StatusInternalServerError
			errorPayload = unexpectedError
		}

		err = sendJSON(w, statusCode, errorPayload)
		if err != nil {
			h.logger.Error("failed to send error json response", zap.Error(err))
		}
		return
	}

	err = sendJSON(w, http.StatusCreated, StudentRegisterResponse{studentID})
	if err != nil {
		h.logger.Error("failed to send json response", zap.Error(err))
	}
}
