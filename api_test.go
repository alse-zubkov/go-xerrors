package xerrors

import (
	"errors"
	"testing"
)

func TestNewCode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		code := NewCode("test::code", Data{"key": "value"})
		if code == nil {
			t.Fatal("expected non-nil code")
		}
		if code.Metadata() == nil {
			t.Fatal("expected non-nil metadata")
		}
	})

	t.Run("nil_key_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil key")
			}
		}()
		NewCode(nil, nil)
	})
}

func TestNew(t *testing.T) {
	code := NewCode("test::code", nil)

	t.Run("success", func(t *testing.T) {
		err := New(code, Data{"key": "value"})
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.Code() != code {
			t.Fatal("unexpected code")
		}
		if err.Type() != TypeSimple {
			t.Fatal("expected TypeSimple")
		}
	})

	t.Run("nil_code_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil code")
			}
		}()
		New(nil, nil)
	})

	t.Run("nil_data_ok", func(t *testing.T) {
		err := New(code, nil)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
	})

	t.Run("data_with_internal_fields_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for internal fields in data")
			}
		}()
		New(code, Data{InternalDataKeyWrappedError: "value"})
	})
}

func TestWrap(t *testing.T) {
	code := NewCode("test::code", nil)
	innerErr := errors.New("inner error")

	t.Run("success", func(t *testing.T) {
		err := Wrap(code, Data{"key": "value"}, innerErr)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.Code() != code {
			t.Fatal("unexpected code")
		}
		if err.Type() != TypeWrapper {
			t.Fatal("expected TypeWrapper")
		}
	})

	t.Run("unwrap_returns_inner_error", func(t *testing.T) {
		err := Wrap(code, nil, innerErr)
		if errors.Unwrap(err) != innerErr {
			t.Fatal("expected Unwrap() to return inner error")
		}
	})

	t.Run("data_excludes_internal_key", func(t *testing.T) {
		err := Wrap(code, Data{"key": "value"}, innerErr)
		if _, ok := err.Data()[InternalDataKeyWrappedError]; ok {
			t.Fatal("Data() should not contain internal key")
		}
	})

	t.Run("nested_returns_nil", func(t *testing.T) {
		err := Wrap(code, nil, innerErr)
		if err.Nested() != nil {
			t.Fatal("expected Nested() to return nil for wrapper")
		}
	})

	t.Run("unwrapped_simple_error_returns_nil", func(t *testing.T) {
		simpleErr := New(code, nil)
		if simpleErr.Unwrap() != nil {
			t.Fatal("expected Unwrap() to return nil for simple error")
		}
	})

	t.Run("nil_code_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil code")
			}
		}()
		Wrap(nil, nil, innerErr)
	})

	t.Run("nil_err_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil err")
			}
		}()
		Wrap(code, nil, nil)
	})

	t.Run("data_with_internal_fields_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for internal fields in data")
			}
		}()
		Wrap(code, Data{InternalDataKeyAggregatedErrors: "value"}, innerErr)
	})
}

func TestAggregate(t *testing.T) {
	code := NewCode("test::code", nil)
	err1 := New(code, Data{"err": "1"})
	err2 := New(code, Data{"err": "2"})

	t.Run("success", func(t *testing.T) {
		err := Aggregate(code, Data{"key": "value"}, err1, err2)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.Code() != code {
			t.Fatal("unexpected code")
		}
		if err.Type() != TypeAggregator {
			t.Fatal("expected TypeAggregator")
		}
	})

	t.Run("nested_returns_errors", func(t *testing.T) {
		err := Aggregate(code, nil, err1, err2)
		nested := err.Nested()
		if len(nested) != 2 {
			t.Fatalf("expected 2 nested errors, got %d", len(nested))
		}
		if nested[0] != err1 || nested[1] != err2 {
			t.Fatal("unexpected nested errors")
		}
	})

	t.Run("data_excludes_internal_key", func(t *testing.T) {
		err := Aggregate(code, Data{"key": "value"}, err1)
		if _, ok := err.Data()[InternalDataKeyAggregatedErrors]; ok {
			t.Fatal("Data() should not contain internal key")
		}
	})

	t.Run("unwrap_returns_nil", func(t *testing.T) {
		err := Aggregate(code, nil, err1, err2)
		if err.Unwrap() != nil {
			t.Fatal("expected Unwrap() to return nil for aggregator")
		}
	})

	t.Run("nil_code_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil code")
			}
		}()
		Aggregate(nil, nil, err1)
	})

	t.Run("nil_err_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil err")
			}
		}()
		Aggregate(code, nil, nil)
	})

	t.Run("empty_errs_ok", func(t *testing.T) {
		err := Aggregate(code, nil)
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if len(err.Nested()) != 0 {
			t.Fatal("expected empty nested errors")
		}
	})

	t.Run("data_with_internal_fields_panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for internal fields in data")
			}
		}()
		Aggregate(code, Data{InternalDataKeyWrappedError: "value"}, err1)
	})
}
