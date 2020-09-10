package backapi

import "github.com/pkg/errors"

type ErrorWithKey struct {
	Key string
}

func (e *ErrorWithKey) Error() string {
	return e.Key
}

// NotFoundError

// NewNotFoundError is used when we cannot find a specified resource
func NewNotFoundError(key string) error {
	return errors.WithStack(&NotFoundError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// NotFoundError is used when we cannot find a specified resource
type NotFoundError struct {
	ErrorWithKey
}

// IsNotFoundError identify an error as an NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := errors.Cause(err).(*NotFoundError)

	return ok
}

// InternalServerError

// NewInternalServerError is used when an error unexpected appears
func NewInternalServerError(key string) error {
	return errors.WithStack(&InternalServerError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// InternalServerError is used when an error unexpected appears
type InternalServerError struct {
	ErrorWithKey
}

// IsInternalServerError identify an error as an InternalServerError
func IsInternalServerError(err error) bool {
	_, ok := errors.Cause(err).(*InternalServerError)

	return ok
}

// UnauthorizedError

// NewUnauthorizedError is used when the action is not authorized
func NewUnauthorizedError(key string, subjectAndMessage ...string) error {
	return errors.WithStack(&UnauthorizedError{
		ErrorWithKey: ErrorWithKey{key},
	})
}

// UnauthorizedError is used when action is not authorized
type UnauthorizedError struct {
	ErrorWithKey
}

// IsUnauthorizedError identifies an error as UnauthorizedError
func IsUnauthorizedError(err error) bool {
	_, ok := errors.Cause(err).(*UnauthorizedError)

	return ok
}

// ResourceAlreadyCreatedError

// NewResourceAlreadyExist is used when a resource already exist and could not be created another time
func NewResourceAlreadyCreatedError(key string) error {
	return errors.WithStack(&ResourceAlreadyCreatedError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// ResourceAlreadyCreatedError is used when a resource already exist and could not be created another time
type ResourceAlreadyCreatedError struct {
	ErrorWithKey
}

// IsResourceAlreadyCreatedError identify an error as a ResourceAlreadyCreatedError
func IsResourceAlreadyCreatedError(err error) bool {
	_, ok := errors.Cause(err).(*ResourceAlreadyCreatedError)

	return ok
}

// OutdatedResourceError

// NewOutdatedResourceError is used when a resource already exist and could not be created another time
func NewOutdatedResourceError(key string) error {
	return errors.WithStack(&OutdatedResourceError{
		ErrorWithKey: ErrorWithKey{
			Key: key,
		},
	})
}

// ResourceAlreadyCreatedError is used when a resource already exist and could not be created another time
type OutdatedResourceError struct {
	ErrorWithKey
}

// IsResourceAlreadyCreatedError identify an error as a ResourceAlreadyCreatedError
func IsOutdatedResourceError(err error) bool {
	_, ok := errors.Cause(err).(*OutdatedResourceError)

	return ok
}
