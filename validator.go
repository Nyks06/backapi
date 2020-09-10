package backapi

// Validator defines the bahaviour that should be implemented by any oject that can be used to validate payload.
type Validator interface {
	Validate(payload interface{}) error
}
