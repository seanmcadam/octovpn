package errors

import (
	"fmt"

	log "github.com/seanmcadam/loggy"
)

func ErrorNilMethodPointer() error {
	return fmt.Errorf("%s Nil Method Pointer", log.FileLineFunc(2))
}

func ErrorShouldNotBeHere() error {
	return fmt.Errorf("%s Should not reach this code", log.FileLineFunc(2))
}
