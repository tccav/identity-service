package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
	"github.com/tccav/identity-service/pkg/domain/identities/idmocks"
	"github.com/tccav/identity-service/pkg/gateways/httpserver/hsfixtures"
)

func TestAuthenticationHandler_AuthenticateStudent(t *testing.T) {
	t.Parallel()

	validToken := entities.Token{
		ID:             uuid.NewString(),
		UserID:         "12345678910",
		ExpirationDate: time.Now().Add(600 * time.Second),
		Hash:           "jwt_token",
	}

	tt := []struct {
		name             string
		requestBody      string
		expectedUC       entities.Token
		expectedUCErr    error
		expectedResponse any
		expectedStatus   any
	}{
		{
			name:           "should successfully authenticate student",
			requestBody:    hsfixtures.ValidStudentLoginRequestBody,
			expectedUC:     validToken,
			expectedStatus: http.StatusCreated,
			expectedResponse: AuthenticateStudentResponse{
				TokenID:   validToken.ID,
				ExpiresAt: validToken.ExpirationDate.Format(time.RFC3339),
				Token:     validToken.Hash,
			},
		},
		{
			name:             "should fail and receive invalid json response",
			requestBody:      hsfixtures.InvalidJSONRequestBody,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: invalidJSON,
		},
		{
			name:             "should fail because student id is empty",
			requestBody:      `{"secret": "123467"}`,
			expectedUCErr:    identities.ErrEmptyStudentID,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: emptyStudentID,
		},
		{
			name:             "should fail because secret is empty",
			requestBody:      `{"student_id": "1234678910"}`,
			expectedUCErr:    identities.ErrEmptySecret,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: emptySecret,
		},
		{
			name:             "should fail because student is not registered",
			requestBody:      hsfixtures.ValidStudentLoginRequestBody,
			expectedUCErr:    identities.ErrStudentNotFound,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: invalidCredentials,
		},
		{
			name:             "should fail because secret doesn't match",
			requestBody:      hsfixtures.ValidStudentLoginRequestBody,
			expectedUCErr:    identities.ErrSecretsDontMatch,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: invalidCredentials,
		},
		{
			name:             "should fail because an unexpected error happened",
			requestBody:      hsfixtures.ValidStudentLoginRequestBody,
			expectedUCErr:    errors.New("unexpected error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: unexpectedError,
		},
	}
	for _, testCase := range tt {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// prepare
			logger := zap.NewNop()

			useCase := idmocks.AuthenticationUseCasesMock{
				AuthenticateStudentFunc: func(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error) {
					return tc.expectedUC, tc.expectedUCErr
				},
			}

			expectedResponse, err := json.Marshal(tc.expectedResponse)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/v1/identities/students/login",
				bytes.NewReader([]byte(tc.requestBody)))

			h := NewAuthenticationHandler(logger, &useCase)

			// test
			h.AuthenticateStudent(w, r)

			// assert
			assert.Equal(t, string(expectedResponse), strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
