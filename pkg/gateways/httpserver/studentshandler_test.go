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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
	"github.com/tccav/identity-service/pkg/domain/identities/idmocks"
	"github.com/tccav/identity-service/pkg/gateways/httpserver/hsfixtures"
)

func TestStudentsHandler_RegisterStudent(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name             string
		requestBody      string
		expectedUC       string
		expectedUCErr    error
		expectedResponse any
		expectedStatus   any
	}{
		{
			name:             "should successfully create student",
			requestBody:      hsfixtures.ValidStudentRequestBody,
			expectedUC:       "123451271",
			expectedResponse: StudentRegisterResponse{ID: "123451271"},
			expectedStatus:   http.StatusCreated,
		},
		{
			name:             "should fail due to invalid json",
			requestBody:      hsfixtures.InvalidJSONRequestBody,
			expectedResponse: invalidJSON,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to invalid student id",
			requestBody:      hsfixtures.InvalidIDRequestBody,
			expectedUCErr:    entities.ErrInvalidStudentID,
			expectedResponse: invalidStudentID,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to invalid course id",
			requestBody:      hsfixtures.InvalidCourseIDRequestBody,
			expectedUCErr:    identities.ErrInvalidCourseID,
			expectedResponse: invalidCourseID,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to invalid cpf",
			requestBody:      hsfixtures.InvalidCPFRequestBody,
			expectedUCErr:    entities.ErrInvalidCPF,
			expectedResponse: invalidCPF,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to invalid email",
			requestBody:      hsfixtures.InvalidEmailRequestBody,
			expectedUCErr:    entities.ErrInvalidEmail,
			expectedResponse: invalidEmail,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to invalid birth date",
			requestBody:      hsfixtures.InvalidBirthDateRequestBody,
			expectedUCErr:    entities.ErrInvalidBirthDate,
			expectedResponse: invalidBirthDate,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "should fail due to unexpected error from use case",
			requestBody:      hsfixtures.ValidStudentRequestBody,
			expectedUCErr:    errors.New("unexpected"),
			expectedResponse: unexpectedError,
			expectedStatus:   http.StatusInternalServerError,
		},
	}
	for _, testCase := range tt {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// prepare
			logger := zap.NewNop()
			registerUseCasesMock := idmocks.RegisterUseCasesMock{RegisterStudentFunc: func(ctx context.Context, input identities.RegisterStudentInput) (string, error) {
				return tc.expectedUC, tc.expectedUCErr
			}}

			expectedResponse, err := json.Marshal(tc.expectedResponse)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/v1/identities/students",
				bytes.NewReader([]byte(tc.requestBody)))

			h := NewStudentsHandler(&registerUseCasesMock, logger)

			// test
			h.RegisterStudent(w, r)

			// assert
			assert.Equal(t, string(expectedResponse), strings.TrimSpace(w.Body.String()))
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
