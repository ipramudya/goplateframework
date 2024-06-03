package errs

import (
	"fmt"
	"runtime"
)

type ErrCode struct {
	value int
}

func (ec ErrCode) Value() int {
	return ec.value
}

func (ec ErrCode) String() string {
	return codeNames[ec]
}

func (ec ErrCode) HTTP() int {
	return httpStatus[ec]
}

type Error struct {
	Code     ErrCode `json:"-"`
	CodeName string  `json:"code_name"`
	Message  string  `json:"message"`
	FuncName string  `json:"-"`
	FileName string  `json:"-"`
}

func New(code ErrCode, err error) *Error {
	pc, filename, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		CodeName: code.String(),
		Message:  err.Error(),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func Newf(code ErrCode, format string, v ...any) *Error {
	pc, filename, line, _ := runtime.Caller(1)

	return &Error{
		Code:     code,
		CodeName: code.String(),
		Message:  fmt.Sprintf(format, v...),
		FuncName: runtime.FuncForPC(pc).Name(),
		FileName: fmt.Sprintf("%s:%d", filename, line),
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) HTTPStatus() int {
	return httpStatus[e.Code]
}

func (e *Error) Debug() string {
	return fmt.Sprintf("code: %s, message: %s, func: %s, file: %s", e.Code, e.Message, e.FuncName, e.FileName)
}
