package internal

import (
	"errors"
	"fmt"
)

var ErrInvalidRequest = errors.New("invalid request")

type Error struct {
	err error
	msg string
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}
	return fmt.Sprintf("%s: %s", e.msg, e.err.Error())
}

func WrapErr(err error, msg string) *Error {
	return &Error{
		err: err,
		msg: msg,
	}
}
