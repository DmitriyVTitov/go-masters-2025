package interfaces

import (
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

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

type T struct{}

func (t *T) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (t *T) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Базовая механика работы с интерфейсами.
func basics() {
	// Конкретный тип данных.
	var t *T

	// Интерфейс.
	var reader Reader

	// Типичный вопрос на собеседованиях про равенство интерфейсной переменной и nil.
	if reader == nil {
		log.Info().Msg("r is nil")
	} else {
		log.Info().Msg("r is not nil")
	}

	if t == nil {
		log.Info().Msg("t is nil")
	} else {
		log.Info().Msg("t is not nil")
	}

	reader = t
	log.Info().Msg("r = t")

	if reader == nil {
		log.Info().Msg("r is nil")
	} else {
		log.Info().Msg("r is not nil")
	}

	reader.Read(nil)

	// Проверка, что тип данных имплементирует интерфейс
	// (используется в json.Marshal/Unmarshal).
	var a any = t
	if _, ok := a.(Reader); ok {
		log.Info().Msgf("%T is Reader", a)
	} else {
		log.Info().Msgf("%T is not Reader", a)
	}

	// Type assertion (interface).
	rw, ok := reader.(ReadWriter)
	if ok {
		log.Info().Msg("r is ReadWriter")
		rw.Write(nil)
	}

	// Type assertion (type).
	val, ok := rw.(*T)
	if ok {
		log.Info().Msg("rw is *T")
		val.Read(nil)
		val.Write(nil)
	}

	// Type switch.
	log.Info().Msg("Type switch")
	switch v := rw.(type) {
	case nil:
		log.Info().Msg("rw is nil")
	case Reader:
		log.Info().Msg("rw is Reader")
		v.Read(nil)
	case Writer:
		log.Info().Msg("rw is Writer")
	case ReadWriter:
		log.Info().Msg("rw is ReadWriter")
	default:
		log.Info().Msg("rw is unknown type")
	}
}
