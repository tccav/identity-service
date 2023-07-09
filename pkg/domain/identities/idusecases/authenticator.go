package idusecases

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type tokenMaker interface {
	createToken(ctx context.Context, userID string) (entities.Token, error)
	verifyToken(ctx context.Context, hash string) error
}

type Config interface {
	TokenSecret() string
	TokenIssuer() string
	TokenDuration() time.Duration
}

type StudentAuthenticator struct {
	tokenMaker
	studentsRepository identities.StudentListerRepository
	tracer             trace.Tracer
}

func NewStudentJWTAuthenticator(studentRepository identities.StudentListerRepository, tokenRepository identities.TokenRegistererRepository, config Config) StudentAuthenticator {
	tracer := otel.Tracer(tracerName)

	maker := jwtTokenMaker{
		secret:     config.TokenSecret(),
		issuer:     config.TokenIssuer(),
		duration:   config.TokenDuration(),
		repository: tokenRepository,
		tracer:     tracer,
	}

	return StudentAuthenticator{
		tokenMaker:         maker,
		studentsRepository: studentRepository,
		tracer:             tracer,
	}
}

func (s StudentAuthenticator) AuthenticateStudent(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error) {
	ctx, span := s.tracer.Start(ctx, "StudentAuthenticator.AuthenticateStudent")
	defer span.End()
	if input.StudentID == "" {
		span.RecordError(identities.ErrEmptyStudentID)
		return entities.Token{}, identities.ErrEmptyStudentID
	}

	if input.StudentSecret == "" {
		span.RecordError(identities.ErrEmptySecret)
		return entities.Token{}, identities.ErrEmptySecret
	}

	registeredSecret, err := s.studentsRepository.GetStudentSecret(ctx, input.StudentID)
	if err != nil {
		span.RecordError(err)
		return entities.Token{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(registeredSecret), []byte(input.StudentSecret))
	if err != nil {
		span.RecordError(err)
		return entities.Token{}, err
	}

	token, err := s.createToken(ctx, input.StudentID)
	if err != nil {
		span.RecordError(err)
		return entities.Token{}, err
	}

	return token, nil
}

func (s StudentAuthenticator) VerifyAuth(ctx context.Context, hash string) error {
	ctx, span := s.tracer.Start(ctx, "StudentAuthenticator.VerifyAuth")
	defer span.End()

	if hash == "" {
		span.RecordError(identities.ErrEmptyToken)
		return identities.ErrEmptyToken
	}

	err := s.tokenMaker.verifyToken(ctx, hash)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
