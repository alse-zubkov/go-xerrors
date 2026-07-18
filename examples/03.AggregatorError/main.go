package main

import (
	"fmt"

	"github.com/alse-zubkov/go-xerrors"
)

var CodeFieldRequired = xerrors.NewCode("validation::field-required", xerrors.Data{})
var CodeFieldWrongValue = xerrors.NewCode("validation::field-wrong-length", xerrors.Data{})

var CodeManifestInvalid = xerrors.NewCode("validation::manifest-invalid", xerrors.Data{
	"description": "Manifest is invalid, unable perform the operation",
	"http.code": 400,
})

func ValidateManifest() error {
	xerrs := []xerrors.Error{
		xerrors.New(CodeFieldRequired, xerrors.Data{"field": "name"}),
		xerrors.New(CodeFieldWrongValue, xerrors.Data{"field": "id", "expected-format":"uuid"}),
	}
	return xerrors.Aggregate(
		CodeManifestInvalid,
		xerrors.Data{},
		xerrs...,
	)
}

func main() {
	err := ValidateManifest()
	if xerr, ok := err.(xerrors.Error); ok {
		fmt.Println(xerr)
		fmt.Println("Code:", xerr.Code().Key())
		fmt.Println("Data:", xerr.Data())
		fmt.Println("Type:", xerr.Type())
		fmt.Println("Errors:", xerr.Nested())
	}
}
