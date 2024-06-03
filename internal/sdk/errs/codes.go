package errs

import (
	"net/http"
)

var (
	OK                 = ErrCode{value: 0}
	NoContent          = ErrCode{value: 1}
	Canceled           = ErrCode{value: 2}
	Unknown            = ErrCode{value: 3}
	InvalidArgument    = ErrCode{value: 4}
	DeadlineExceeded   = ErrCode{value: 5}
	NotFound           = ErrCode{value: 6}
	AlreadyExists      = ErrCode{value: 7}
	PermissionDenied   = ErrCode{value: 8}
	ResourceExhausted  = ErrCode{value: 9}
	FailedPrecondition = ErrCode{value: 10}
	Aborted            = ErrCode{value: 11}
	OutOfRange         = ErrCode{value: 12}
	Unimplemented      = ErrCode{value: 13}
	Internal           = ErrCode{value: 14}
	Unavailable        = ErrCode{value: 15}
	DataLoss           = ErrCode{value: 16}
	Unauthenticated    = ErrCode{value: 17}
	InvalidCredentials = ErrCode{value: 18}
)

var codeNames = map[ErrCode]string{
	OK:                 "ok",
	NoContent:          "ok_no_content",
	Canceled:           "canceled",
	Unknown:            "unknown",
	InvalidArgument:    "invalid_argument",
	DeadlineExceeded:   "deadline_exceeded",
	NotFound:           "not_found",
	AlreadyExists:      "already_exists",
	PermissionDenied:   "permission_denied",
	ResourceExhausted:  "resource_exhausted",
	FailedPrecondition: "failed_precondition",
	Aborted:            "aborted",
	OutOfRange:         "out_of_range",
	Unimplemented:      "unimplemented",
	Internal:           "internal",
	Unavailable:        "unavailable",
	DataLoss:           "data_loss",
	Unauthenticated:    "unauthenticated",
	InvalidCredentials: "invalid_credentials",
}

var httpStatus = map[ErrCode]int{
	OK:                 http.StatusOK,
	NoContent:          http.StatusNoContent,
	Canceled:           http.StatusGatewayTimeout,
	Unknown:            http.StatusInternalServerError,
	InvalidArgument:    http.StatusBadRequest,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	NotFound:           http.StatusNotFound,
	AlreadyExists:      http.StatusConflict,
	PermissionDenied:   http.StatusForbidden,
	ResourceExhausted:  http.StatusTooManyRequests,
	FailedPrecondition: http.StatusBadRequest,
	Aborted:            http.StatusConflict,
	OutOfRange:         http.StatusBadRequest,
	Unimplemented:      http.StatusNotImplemented,
	Internal:           http.StatusInternalServerError,
	Unavailable:        http.StatusServiceUnavailable,
	DataLoss:           http.StatusInternalServerError,
	Unauthenticated:    http.StatusUnauthorized,
	InvalidCredentials: http.StatusUnauthorized,
}
