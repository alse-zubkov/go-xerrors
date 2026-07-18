package xerrors

import "fmt"

// Data is a map of string keys to any values used for error metadata.
type Data map[string]any

// CodeKey is an identifier for error codes.
type CodeKey any

// Code represents an error code with a key and metadata.
type Code interface {
	Key() CodeKey
	Metadata() Data
}

// Type represents the type of an error.
type Type string

const (
	// TypeSimple indicates a simple error.
	TypeSimple Type = "simple"
	// TypeWrapper indicates a wrapper error that contains a source error.
	TypeWrapper Type = "wrapper"
	// TypeAggregator indicates an aggregator error that contains multiple errors.
	TypeAggregator Type = "aggregator"
)

// Error represents an error with code, type, and data.
type Error interface {
	error
	// Code returns error code
	Code() Code
	// Type returns error type
	Type() Type
	// Data returns error attributes
	Data() Data
	// Unwrap returns source error in case of TypeWrapper
	Unwrap() error
	// Nested returns a slice of aggregated Error(s)
	Nested() []Error
}

// codeBase is the base implementation of Code.
type codeBase struct {
	key      CodeKey
	metadata Data
}

// Key returns the code key.
func (c *codeBase) Key() CodeKey {
	return c.key
}

// Metadata returns the code metadata.
func (c *codeBase) Metadata() Data {
	return c.metadata
}

type errorBase struct {
	code Code
    t    Type
	data Data
}
// Error implements the error interface.
func (e *errorBase) Error() string {
	if e.code == nil {
		return "Xerror{nil}"
	}
	if len(e.data) < 1 {
		return fmt.Sprintf("Xerror{%v}", e.code.Key())
	}
	userData := e.Data()
	if userData == nil {
		return fmt.Sprintf("Xerror{%v}", e.code.Key())
	}
	return fmt.Sprintf("Xerror{%v;%v}", e.code.Key(), userData)
}

// Code returns the error code.
func (e *errorBase) Code() Code {
	return e.code
}

// Type returns the error type.
func (e *errorBase) Type() Type {
	return e.t
}

// Data returns the error data without internal fields.
func (e *errorBase) Data() Data {
	if e.data == nil {
		return nil
	}
	result := make(Data, len(e.data))
	for k, v := range e.data {
		if k == InternalDataKeyWrappedError || k == InternalDataKeyAggregatedErrors {
			continue
		}
		result[k] = v
	}
	return result
}

// Unwrap returns the wrapped error for TypeWrapper errors.
func (e *errorBase) Unwrap() error {
	if e.t != TypeWrapper {
		return nil
	}
	if e.data == nil {
		return nil
	}
	if err, ok := e.data[InternalDataKeyWrappedError]; ok {
		return err.(error)
	}
	return nil
}

// Nested returns the nested errors for TypeAggregator errors.
func (e *errorBase) Nested() []Error {
	if e.t != TypeAggregator {
		return nil
	}
	if e.data == nil {
		return nil
	}
	if errs, ok := e.data[InternalDataKeyAggregatedErrors]; ok {
		return errs.([]Error)
	}
	return nil
}
