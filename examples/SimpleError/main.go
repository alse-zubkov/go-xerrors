package main

import (
	"fmt"

	"github.com/alse-zubkov/go-xerrors"
)

var CodeEntityNotFound = xerrors.NewCode("storage::entity-not-found", xerrors.Data{
	"description": "Entity was not found",
	"http.code":   404,
})

func FindUserByID(id any) error {
	return xerrors.New(CodeEntityNotFound,
		xerrors.Data{
			"entity-type": "user",
			"entity-id":   id,
		},
	)
}

func main() {
	err := FindUserByID(123)
	if xerr, ok := err.(xerrors.Error); ok {
		fmt.Println(xerr)
		fmt.Println("Code:", xerr.Code().Key())
		fmt.Println("Data:", xerr.Data())
		fmt.Println("Type:", xerr.Type())
	}
}
