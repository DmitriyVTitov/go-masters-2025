package errs

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// ***
// Все типы ошибок приложения находятся в отдельном пакете.
// ***

type BaseError struct {
	Message string
	Time    time.Time
}

func (e *BaseError) Error() string {
	return e.Message
}

func NewBaseError(message string) *BaseError {
	return &BaseError{
		Message: message,
		Time:    time.Now(),
	}
}

type ErrBadRequest struct {
	BaseError
}

func NewErrBadRequest(message string) *ErrBadRequest {
	return &ErrBadRequest{
		BaseError: BaseError{
			Message: message,
			Time:    time.Now(),
		},
	}

}

type ErrUnauthorized struct {
	BaseError
}

func NewErrUnauthorized(message string) *ErrUnauthorized {
	return &ErrUnauthorized{
		BaseError: BaseError{
			Message: message,
			Time:    time.Now(),
		},
	}
}

// ***
// Если приложение предоставляет HTTP API, то можно определить
// функцию, которая выводит пользователю сообщение об ошибке,
// в соответствии с типами ошибок из пакета
// ***

type API struct{}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// WriteError отправляет пользователю код ошибки и соответствующе сообщение
func (api *API) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	resp := ErrorResponse{
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")

	// Код ответа в соответствии с типом ошибки
	switch err.(type) {
	case *ErrBadRequest:
		resp.StatusCode = http.StatusBadRequest
	case *ErrUnauthorized:
		resp.StatusCode = http.StatusUnauthorized
	default:
		resp.StatusCode = http.StatusInternalServerError
	}

	// Альтернативный вариант проверки типа ошибки
	var errBadRequest *ErrBadRequest
	if errors.As(err, &errBadRequest) {
		resp.StatusCode = http.StatusBadRequest
	}

	var errUnauthorized *ErrUnauthorized
	if errors.As(err, &errUnauthorized) {
		resp.StatusCode = http.StatusUnauthorized
	}

	// Установка кода ответа HTTP
	w.WriteHeader(resp.StatusCode)

	// Отправка пользователю JSON с ошибкой
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Err(err).Send()
	}
}
