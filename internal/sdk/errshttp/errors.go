package errshttp

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorResponse struct {
	Err struct {
		ErrorCode ErrorCode     `json:"-"`
		Code      string        `json:"code"`
		Message   string        `json:"message"`
		RequestID string        `json:"request_id"`
		Details   []ErrorDetail `json:"details,omitempty"`

		// debugger
		FuncName string `json:"-"`
		FileName string `json:"-"`
	} `json:"error"`
}

type ErrorDetail struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func New(code ErrorCode, message string) *ErrorResponse {
	err := new(ErrorResponse)

	pc, filename, line, _ := runtime.Caller(1)

	err.Err.ErrorCode = code
	err.Err.Code = code.String()
	err.Err.Message = message
	err.Err.FuncName = runtime.FuncForPC(pc).Name()
	err.Err.FileName = fmt.Sprintf("%s:%d", filename, line)

	return err
}

func (e *ErrorResponse) AddDetail(detail string) {
	s := strings.Split(detail, ":")

	if len(s) == 2 {
		e.Err.Details = append(
			e.Err.Details,
			ErrorDetail{Field: strings.TrimSpace(s[0]), Reason: strings.TrimSpace(s[1])},
		)
	} else {
		e.Err.Details = append(
			e.Err.Details,
			ErrorDetail{Field: "-", Reason: strings.TrimSpace(detail)},
		)
	}
}

func (err *ErrorResponse) AddRequestID(requestID string) {
	err.Err.RequestID = requestID
}

func (err *ErrorResponse) ErrorForLoggingDebug() string {
	return fmt.Sprintf("[code: %s, message: %s, func: %s, file: %s]", err.Err.Code, err.Err.Message, err.Err.FuncName, err.Err.FileName)
}

func (e *ErrorResponse) Error() string {
	return e.Err.Message
}

func (e *ErrorResponse) HTTPStatus() int {
	return e.Err.ErrorCode.HTTP()
}
