// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package idmocks

import (
	"context"
	"github.com/tccav/identity-service/pkg/domain/entities"
	"github.com/tccav/identity-service/pkg/domain/identities"
	"sync"
)

// Ensure, that RegisterUseCasesMock does implement identities.RegisterUseCases.
// If this is not the case, regenerate this file with moq.
var _ identities.RegisterUseCases = &RegisterUseCasesMock{}

// RegisterUseCasesMock is a mock implementation of identities.RegisterUseCases.
//
//	func TestSomethingThatUsesRegisterUseCases(t *testing.T) {
//
//		// make and configure a mocked identities.RegisterUseCases
//		mockedRegisterUseCases := &RegisterUseCasesMock{
//			RegisterStudentFunc: func(ctx context.Context, input identities.RegisterStudentInput) (string, error) {
//				panic("mock out the RegisterStudent method")
//			},
//		}
//
//		// use mockedRegisterUseCases in code that requires identities.RegisterUseCases
//		// and then make assertions.
//
//	}
type RegisterUseCasesMock struct {
	// RegisterStudentFunc mocks the RegisterStudent method.
	RegisterStudentFunc func(ctx context.Context, input identities.RegisterStudentInput) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// RegisterStudent holds details about calls to the RegisterStudent method.
		RegisterStudent []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input identities.RegisterStudentInput
		}
	}
	lockRegisterStudent sync.RWMutex
}

// RegisterStudent calls RegisterStudentFunc.
func (mock *RegisterUseCasesMock) RegisterStudent(ctx context.Context, input identities.RegisterStudentInput) (string, error) {
	if mock.RegisterStudentFunc == nil {
		panic("RegisterUseCasesMock.RegisterStudentFunc: method is nil but RegisterUseCases.RegisterStudent was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Input identities.RegisterStudentInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockRegisterStudent.Lock()
	mock.calls.RegisterStudent = append(mock.calls.RegisterStudent, callInfo)
	mock.lockRegisterStudent.Unlock()
	return mock.RegisterStudentFunc(ctx, input)
}

// RegisterStudentCalls gets all the calls that were made to RegisterStudent.
// Check the length with:
//
//	len(mockedRegisterUseCases.RegisterStudentCalls())
func (mock *RegisterUseCasesMock) RegisterStudentCalls() []struct {
	Ctx   context.Context
	Input identities.RegisterStudentInput
} {
	var calls []struct {
		Ctx   context.Context
		Input identities.RegisterStudentInput
	}
	mock.lockRegisterStudent.RLock()
	calls = mock.calls.RegisterStudent
	mock.lockRegisterStudent.RUnlock()
	return calls
}

// Ensure, that AuthenticationUseCasesMock does implement identities.AuthenticationUseCases.
// If this is not the case, regenerate this file with moq.
var _ identities.AuthenticationUseCases = &AuthenticationUseCasesMock{}

// AuthenticationUseCasesMock is a mock implementation of identities.AuthenticationUseCases.
//
//	func TestSomethingThatUsesAuthenticationUseCases(t *testing.T) {
//
//		// make and configure a mocked identities.AuthenticationUseCases
//		mockedAuthenticationUseCases := &AuthenticationUseCasesMock{
//			AuthenticateStudentFunc: func(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error) {
//				panic("mock out the AuthenticateStudent method")
//			},
//			VerifyAuthFunc: func(ctx context.Context, hash string) error {
//				panic("mock out the VerifyAuth method")
//			},
//		}
//
//		// use mockedAuthenticationUseCases in code that requires identities.AuthenticationUseCases
//		// and then make assertions.
//
//	}
type AuthenticationUseCasesMock struct {
	// AuthenticateStudentFunc mocks the AuthenticateStudent method.
	AuthenticateStudentFunc func(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error)

	// VerifyAuthFunc mocks the VerifyAuth method.
	VerifyAuthFunc func(ctx context.Context, hash string) error

	// calls tracks calls to the methods.
	calls struct {
		// AuthenticateStudent holds details about calls to the AuthenticateStudent method.
		AuthenticateStudent []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Input is the input argument value.
			Input identities.AuthenticateStudentInput
		}
		// VerifyAuth holds details about calls to the VerifyAuth method.
		VerifyAuth []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Hash is the hash argument value.
			Hash string
		}
	}
	lockAuthenticateStudent sync.RWMutex
	lockVerifyAuth          sync.RWMutex
}

// AuthenticateStudent calls AuthenticateStudentFunc.
func (mock *AuthenticationUseCasesMock) AuthenticateStudent(ctx context.Context, input identities.AuthenticateStudentInput) (entities.Token, error) {
	if mock.AuthenticateStudentFunc == nil {
		panic("AuthenticationUseCasesMock.AuthenticateStudentFunc: method is nil but AuthenticationUseCases.AuthenticateStudent was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Input identities.AuthenticateStudentInput
	}{
		Ctx:   ctx,
		Input: input,
	}
	mock.lockAuthenticateStudent.Lock()
	mock.calls.AuthenticateStudent = append(mock.calls.AuthenticateStudent, callInfo)
	mock.lockAuthenticateStudent.Unlock()
	return mock.AuthenticateStudentFunc(ctx, input)
}

// AuthenticateStudentCalls gets all the calls that were made to AuthenticateStudent.
// Check the length with:
//
//	len(mockedAuthenticationUseCases.AuthenticateStudentCalls())
func (mock *AuthenticationUseCasesMock) AuthenticateStudentCalls() []struct {
	Ctx   context.Context
	Input identities.AuthenticateStudentInput
} {
	var calls []struct {
		Ctx   context.Context
		Input identities.AuthenticateStudentInput
	}
	mock.lockAuthenticateStudent.RLock()
	calls = mock.calls.AuthenticateStudent
	mock.lockAuthenticateStudent.RUnlock()
	return calls
}

// VerifyAuth calls VerifyAuthFunc.
func (mock *AuthenticationUseCasesMock) VerifyAuth(ctx context.Context, hash string) error {
	if mock.VerifyAuthFunc == nil {
		panic("AuthenticationUseCasesMock.VerifyAuthFunc: method is nil but AuthenticationUseCases.VerifyAuth was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Hash string
	}{
		Ctx:  ctx,
		Hash: hash,
	}
	mock.lockVerifyAuth.Lock()
	mock.calls.VerifyAuth = append(mock.calls.VerifyAuth, callInfo)
	mock.lockVerifyAuth.Unlock()
	return mock.VerifyAuthFunc(ctx, hash)
}

// VerifyAuthCalls gets all the calls that were made to VerifyAuth.
// Check the length with:
//
//	len(mockedAuthenticationUseCases.VerifyAuthCalls())
func (mock *AuthenticationUseCasesMock) VerifyAuthCalls() []struct {
	Ctx  context.Context
	Hash string
} {
	var calls []struct {
		Ctx  context.Context
		Hash string
	}
	mock.lockVerifyAuth.RLock()
	calls = mock.calls.VerifyAuth
	mock.lockVerifyAuth.RUnlock()
	return calls
}
