package main

import (
	"fmt"
	"net/http"

	"github.com/kamva/hexa"
)

type ExampleErrorType struct {
	Name string `json:"name"`
}

var SampleHexaErr = hexa.NewError(http.StatusUnprocessableEntity, "arranger.examples.err", nil)

func (e ExampleErrorType) Error() string {
	return fmt.Sprint("err_msg-> with name: ", e.Name)
}

func NewErr(name string) ExampleErrorType {
	return ExampleErrorType{
		Name: name,
	}
}

var _ error = &ExampleErrorType{}
