package xerrors

import "fmt"

const InternalDataKeyWrappedError = "xerrors::wrapped-error"
const InternalDataKeyAggregatedErrors = "xerrors::aggregated-errors"

var internalValueOperation = "Operation"
var internalValueValues = "Values"
var internalValueMBNN = "must be not nil"

var internalCodeIllegalValue = &codeBase{key: "xerrors::illegal-value", metadata: nil}

// NewCode creates a simple code to be used as base for errors.
func NewCode(key CodeKey, metadata Data) Code {
	valErrs := Data{}
	validateNullValue("key", key, valErrs)
	panicIllegalValueIfNeeded("NewCode", valErrs)

	return &codeBase{key: key, metadata: metadata}
}

// New creates a simple code-based error with data.
func New(code Code, data Data) Error {
	valErrs := Data{}
	validateNullValue("code", code, valErrs)
	validateDataOnSystemFields(data, valErrs)
	panicIllegalValueIfNeeded("New", valErrs)
	return &errorBase{code: code, t: TypeSimple, data: data}
}

// Wrap creates a wrapper around source error with data.
func Wrap(code Code, data Data, err error) Error {
	valErrs := Data{}
	validateNullValue("code", code, valErrs)
	validateNullValue("err", err, valErrs)
	validateDataOnSystemFields(data, valErrs)
	panicIllegalValueIfNeeded("Wrap", valErrs)
	if len(data) < 1 {
		data = Data{}
	}
	data[InternalDataKeyWrappedError] = err
	return &errorBase{code: code, t: TypeWrapper, data: data}
}

// Aggregate creates an aggregator error that contains multiple errors.
func Aggregate(code Code, data Data, errs ...Error) Error {
	valErrs := Data{}
	validateNullValue("code", code, valErrs)
	for i, e := range errs {
		validateNullValue(fmt.Sprintf("errs[%d]", i), e, valErrs)
	}
	validateDataOnSystemFields(data, valErrs)
	panicIllegalValueIfNeeded("Aggregate", valErrs)
	if len(data) < 1 {
		data = Data{}
	}
	data[InternalDataKeyAggregatedErrors] = errs
	return &errorBase{code: code, t: TypeAggregator, data: data}
}

func validateNullValue(arg string, value any, valErrs Data) {
	if value == nil {
		valErrs[arg] = internalValueMBNN
	}
}

func validateDataOnSystemFields(value Data, valErrs Data) {
    if value == nil {
        return
    }
	_, e1 := value[InternalDataKeyWrappedError]
	_, e2 := value[InternalDataKeyAggregatedErrors]
	if e1 || e2 {
		valErrs["data"] = "data contains xerrors internal fields"
	}
}

func panicIllegalValueIfNeeded(operation string, values Data) {
	if len(values) > 0 {
		panic(New(internalCodeIllegalValue, Data{
			internalValueOperation: operation,
			internalValueValues:    values,
		}))
	}
}
