package main

import (
	"errors"
	"fmt"

	"github.com/alse-zubkov/go-xerrors"
)

var CodeEntityNotFound = xerrors.NewCode("storage::entity-not-found", xerrors.Data{
	"description": "Entity was not found",
	"http.code":   404,
})

func FindUserDAO(id any) error {
	return errors.New("no rows")
}

func FindUserByID(id any) error {
	daoError := FindUserDAO(id)
	return xerrors.Wrap(CodeEntityNotFound,
		xerrors.Data{
			"entity-type": "user",
			"entity-id":   id,
		},
		daoError,
	)
}

func main() {
	err := FindUserByID(123)
	if xerr, ok := err.(xerrors.Error); ok {
		fmt.Println(xerr)
		fmt.Println("Code:", xerr.Code().Key())
		fmt.Println("Data:", xerr.Data())
		fmt.Println("Type:", xerr.Type())
		fmt.Println("Error:", xerr.Unwrap())
	}
}
