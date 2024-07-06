package errshttp

import "net/http"

type ErrorCode int

const (
	OK ErrorCode = iota
	NoContent
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
	InvalidCredentials
)

func (ec *ErrorCode) String() string {
	codenames := map[ErrorCode]string{
		OK:                 "OK",
		NoContent:          "NO_CONTENT",
		Canceled:           "CANCELED",
		Unknown:            "UNKNOWN",
		InvalidArgument:    "INVALID_ARGUMENT",
		DeadlineExceeded:   "DEADLINE_EXCEEDED",
		NotFound:           "NOT_FOUND",
		AlreadyExists:      "ALREADY_EXISTS",
		PermissionDenied:   "PERMISSION_DENIED",
		ResourceExhausted:  "RESOURCE_EXHAUSTED",
		FailedPrecondition: "FAILED_PRECONDITION",
		Aborted:            "ABORTED",
		OutOfRange:         "OUT_OF_RANGE",
		Unimplemented:      "UNIMPLEMENTED",
		Internal:           "INTERNAL",
		Unavailable:        "UNAVAILABLE",
		DataLoss:           "DATA_LOSS",
		Unauthenticated:    "UNAUTHENTICATED",
		InvalidCredentials: "INVALID_CREDENTIALS",
	}

	return codenames[*ec]
}

func (ec *ErrorCode) HTTP() int {
	httpstatus := map[ErrorCode]int{
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

	return httpstatus[*ec]
}
