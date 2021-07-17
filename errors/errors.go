package errors

import (
	"fmt"
)

const (
	errorStrFmt = "[Error Code: %d] Reason: %v"
)

type SherlockErrTemplate struct {
	ErrorCode   int    `json:"code"`
	ErrorReason string `json:"reason"`
}

type SherlockErr struct {
	*SherlockErrTemplate
	Err error `json:"error"`
}

func (e *SherlockErr) Error() string {
	return e.stringify()
}

func New(templ *SherlockErrTemplate, err error) *SherlockErr {
	e := &SherlockErr{
		SherlockErrTemplate: templ,
		Err:                 err,
	}
	//e.Error = err
	return e
}

func (e *SherlockErr) stringify() string {
	return fmt.Sprintf(errorStrFmt, e.ErrorCode, e.ErrorReason)
}
