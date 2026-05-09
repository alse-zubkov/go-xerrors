# go-xerrors

Simple library for code-based errors enriched with data and metadata.

## Features

- Create simple errors with custom codes and data
- Wrap underlying errors preserving the original
- Aggregate multiple errors together
- Pure Go implementation, no external dependencies

## Installation

```bash
go get github.com/alse-zubkov/go-xerrors
```

## Quick Start

```go
import "github.com/alse-zubkov/go-xerrors"

// Create error code
code := xerrors.NewCode("domain::example", xerrors.Data{"field": "description"})

// Create simple error
err := xerrors.New(code, xerrors.Data{"key": "value"})
```

## Error Types

### Simple Error

```go
code := xerrors.NewCode("user::not-found", nil)
err := xerrors.New(code, xerrors.Data{"userId": 123})
```

### Wrapped Error

```go
innerErr := errors.New("database connection failed")
err := xerrors.Wrap(code, xerrors.Data{"query": "SELECT *"}, innerErr)

// Access wrapped error
cause := err.Unwrap()
```

### Aggregated Errors

```go
err1 := xerrors.New(code, xerrors.Data{"field": "email"})
err2 := xerrors.New(code, xerrors.Data{"field": "password"})

aggErr := xerrors.Aggregate(code, xerrors.Data{"count": 2}, err1, err2)

// Access nested errors
nested := aggErr.Nested()
```


## Error Types

- `TypeSimple` - simple error created with `New`
- `TypeWrapper` - error that wraps another error, created with `Wrap`
- `TypeAggregator` - error that contains multiple errors, created with `Aggregate`

## Data Fields

`Data()` returns error metadata excluding internal fields used by `Wrap` and `Aggregate`.
