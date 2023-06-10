package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/tccav/identity-service/pkg/domain/identities"
)

type AuthenticateStudentRequest struct {
	StudentID string `json:"student_id" swaggertype:"string" example:"201210204310"`
	Secret    string `json:"secret" swaggertype:"string" example:"celacanto-provoca-maremoto"`
}

type AuthenticateStudentResponse struct {
	TokenID   string `json:"token_id" swaggertype:"string" format:"uuidv4" example:"1f6a4d3a-38c7-43fe-9790-2408fe595c93"`
	ExpiresAt string `json:"expires_at" swaggertype:"string" format:"datetime" example:"2023-10-18T19:32:00.000Z"`
	Token     string `json:"token" swaggertype:"string" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
}

type AuthenticationHandler struct {
	logger *zap.Logger

	useCase identities.AuthenticationUseCases
}

func NewAuthenticationHandler(logger *zap.Logger, useCase identities.AuthenticationUseCases) AuthenticationHandler {
	return AuthenticationHandler{
		logger:  logger,
		useCase: useCase,
	}
}

// AuthenticateStudent ...
// ShowEntity godoc
// @Summary Authenticate a student
// @Tags Auth
// @Param request body AuthenticateStudentRequest true "Student credentials"
// @Accept json
// @Produce json
// @Success 201 {object} AuthenticateStudentResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/identities/students/login [post]
func (h AuthenticationHandler) AuthenticateStudent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var reqBody AuthenticateStudentRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		h.logger.Error("invalid json received", zap.Error(err))
		err = sendJSON(w, http.StatusBadRequest, invalidJSON)
		if err != nil {
			h.logger.Error("failed to send error response", zap.Error(err))
		}
		return
	}

	token, err := h.useCase.AuthenticateStudent(ctx, identities.AuthenticateStudentInput{
		StudentID:     reqBody.StudentID,
		StudentSecret: reqBody.Secret,
	})
	if err != nil {
		h.logger.Error("unable to authenticate user", zap.Error(err))

		var (
			errorPayload HTTPError
			statusCode   int
		)
		switch {
		case errors.Is(err, identities.ErrEmptyStudentID):
			statusCode = http.StatusBadRequest
			errorPayload = emptyStudentID
		case errors.Is(err, identities.ErrEmptySecret):
			statusCode = http.StatusBadRequest
			errorPayload = emptySecret
		case errors.Is(err, identities.ErrStudentNotFound), errors.Is(err, identities.ErrSecretsDontMatch):
			statusCode = http.StatusBadRequest
			errorPayload = invalidCredentials
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

	err = sendJSON(w, http.StatusCreated, AuthenticateStudentResponse{
		TokenID:   token.ID,
		ExpiresAt: token.ExpirationDate.Format(time.RFC3339),
		Token:     token.Hash,
	})
	if err != nil {
		h.logger.Error("failed to send json response", zap.Error(err))
	}
}

// VerifyAuthentication ...
// ShowEntity godoc
// @Summary Verifies if Student Authentication is valid
// @Tags Auth
// @Param authorization header string true "Authorization token"
// @Produce json
// @Success 200
// @Failure 400 {object} HTTPError
// @Failure 401 {object} HTTPError
// @Failure 403 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /v1/identities/students/verify-auth [post]
func (h AuthenticationHandler) VerifyAuthentication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authHeader := strings.Split(strings.TrimSpace(r.Header.Get("authorization")), " ")
	if len(authHeader) != 2 || strings.ToLower(authHeader[0]) != "bearer" {
		err := sendJSON(w, http.StatusForbidden, accessForbidden)
		if err != nil {
			h.logger.Error("failed to send error json response", zap.Error(err))
		}
		return
	}

	token := authHeader[1]

	err := h.useCase.VerifyAuth(ctx, token)
	if err != nil {
		h.logger.Error("unable to verify user auth", zap.Error(err))

		var (
			errorPayload HTTPError
			statusCode   int
		)
		switch {
		case errors.Is(err, identities.ErrTokenNotEmitted), errors.Is(err, identities.ErrMalformedToken):
			statusCode = http.StatusForbidden
			errorPayload = accessForbidden
		case errors.Is(err, identities.ErrTokenExpired):
			statusCode = http.StatusUnauthorized
			errorPayload = accessUnauthorized
			w.Header().Add("WWW-Authenticate", `Bearer realm=".",error="invalid_token",uri="/v1/identities/login"`)
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

	w.WriteHeader(http.StatusOK)
}
