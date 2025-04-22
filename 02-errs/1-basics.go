package errs

import (
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			NoColor:    true,
			TimeFormat: "2006-01-02 15:04:05",
		},
	).With().Timestamp().Logger().With().Caller().Logger()
}

var e error = errors.New("generic error")

// Тип данных myError выполняет контракт интерфейса ощшибки.
type myError struct {
	code    int
	message string
}

func (e *myError) Error() string {
	return e.message
}

var errMyError = &myError{
	code:    10,
	message: "well known error",
}

func caller() {
	err := basics()
	if err != nil {
		// Типичное поведение при ненулевой ошибке - возврат её на уровень выше:
		// if err != nil {
		// 	return defaultValue, err
		// }
		//
		// На верхнем уровне ошибку можно исследовать дополнительно.

		// Проверка на сответствие конкретному экземпляру ошибки.
		if err == errMyError {
			log.Error().Err(err).Msg("error = errMyerror")
		}

		// Проверка на соответствие конкретному типу данных.
		if e, ok := err.(*myError); ok {
			log.Error().Err(e).Msgf("error is myError; code: %v", e.code)
		} else {
			log.Error().Err(err).Msg("error in func basics()")
		}

		// Проверка того, что в цепочке ошибок есть конкретный экземпляр ошибки.
		if errors.Is(err, errMyError) {
			log.Error().Err(err).Msgf("error IS errMyError = true; code: %v", errMyError.code)
		}

		// Проверка того, что в цепочке ошибок есть ошибка определённого типа данных.
		var e *myError
		if errors.As(err, &e) {
			log.Error().Err(e).Msgf("error AS *myError = true; code: %v", e.code)
		}
	}
}

// Ошибка - обычное значение, возвращаемое функцией.
// Ошибка всегда возвращается последней.
func basics() error {
	return fmt.Errorf("basics(): %w", errMyError)
}
