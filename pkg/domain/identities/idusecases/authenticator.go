package idusecases

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
)

type createTokenInput struct {
	userID   string
	duration time.Duration
}

type tokenMaker interface {
	createToken(ctx context.Context, input createTokenInput) (entities.Token, error)
}

type StudentAuthenticator struct {
	tokenMaker
	studentsRepository identities.StudentListerRepository
}

func NewStudentJWTAuthenticator(studentRepository identities.StudentListerRepository, tokenRepository identities.TokenRegistererRepository, secret string) StudentAuthenticator {
	maker := jwtTokenMaker{
		secret:     secret,
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

	token, err := s.createToken(ctx, createTokenInput{
		userID:   input.StudentID,
		duration: defaultTokenDuration,
	})
	if err != nil {
		return entities.Token{}, err
	}

	return token, nil
}
