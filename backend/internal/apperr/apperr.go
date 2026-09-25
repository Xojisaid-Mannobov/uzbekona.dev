// Package apperr ilova xatolarini HTTP status va mashina o'qiy oladigan kod bilan ifodalaydi.
package apperr

import "net/http"

type Error struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (e *Error) Error() string { return e.Message }

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func NotFound(message string) *Error {
	if message == "" {
		message = "Topilmadi"
	}
	return New(http.StatusNotFound, "not_found", message)
}

func BadRequest(message string) *Error {
	return New(http.StatusBadRequest, "bad_request", message)
}

func Unauthorized(message string) *Error {
	if message == "" {
		message = "Avtorizatsiya talab qilinadi"
	}
	return New(http.StatusUnauthorized, "unauthorized", message)
}

func Forbidden(message string) *Error {
	if message == "" {
		message = "Bu amalga ruxsatingiz yo‘q"
	}
	return New(http.StatusForbidden, "forbidden", message)
}

// Validation — maydonlar bo'yicha xatolar (422).
func Validation(fields map[string]string) *Error {
	return &Error{
		Status:  http.StatusUnprocessableEntity,
		Code:    "validation_error",
		Message: "Ma’lumotlarni tekshiring",
		Fields:  fields,
	}
}

// FieldConflict — masalan, slug band bo'lsa (409).
func FieldConflict(field, message string) *Error {
	return &Error{
		Status:  http.StatusConflict,
		Code:    "conflict",
		Message: message,
		Fields:  map[string]string{field: message},
	}
}
