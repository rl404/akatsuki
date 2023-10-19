package validation

// Validator is validating interface.
//
// See usage example in example folder.
type Validator interface {
	// Register custom modifier.
	RegisterModifier(name string, fn func(in string, param ...string) (out string)) error
	// Modify struct field value according to modifier tag.
	// Param `data` should be a pointer.
	Modify(data interface{}) error

	// Register custom validator.
	RegisterValidator(name string, fn func(value interface{}, param ...string) (ok bool)) error
	// Register error message handler.
	RegisterValidatorError(name string, fn func(field string, param ...string) (msg error)) error
	// Validate struct field value according to validator tag.
	// Param `data` should be a pointer.
	Validate(data interface{}) error
}
