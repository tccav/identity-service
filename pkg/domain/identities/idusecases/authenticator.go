package idusecases

import (
	"context"
	"time"

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
}

func NewStudentJWTAuthenticator(studentRepository identities.StudentListerRepository, tokenRepository identities.TokenRegistererRepository, config Config) StudentAuthenticator {
	maker := jwtTokenMaker{
		secret:     config.TokenSecret(),
		issuer:     config.TokenIssuer(),
		duration:   config.TokenDuration(),
		repository: tokenRepository,
	}

	return StudentAuthenticator{
		tokenMaker:         maker,
		studentsRepository: studentRepository,
	}
}

func (s StudentAuthenticator) AuthenticateStudent(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error) {
	if input.StudentID == "" {
		return entities.Token{}, identities.ErrEmptyStudentID
	}

	if input.StudentSecret == "" {
		return entities.Token{}, identities.ErrEmptySecret
	}

	registeredSecret, err := s.studentsRepository.GetStudentSecret(ctx, input.StudentID)
	if err != nil {
		return entities.Token{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(registeredSecret), []byte(input.StudentSecret))
	if err != nil {
		return entities.Token{}, err
	}

	token, err := s.createToken(ctx, input.StudentID)
	if err != nil {
		return entities.Token{}, err
	}

	return token, nil
}

func (s StudentAuthenticator) VerifyAuth(ctx context.Context, hash string) error {
	if hash == "" {
		return identities.ErrEmptyToken
	}
	return s.tokenMaker.verifyToken(ctx, hash)
}
