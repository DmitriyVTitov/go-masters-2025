package interfaces

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

// ***
// Задача на web crawler
//
// Используя тип данных Parser, нужно написать алгоритм обхода веб-страниц
// с уччётом гиперссылок и парсинга ссылок с этих страниц.
// В результате обхода должен быть сформирован массив со всеми найденными ссылками.
type Parser interface {
	Parse(url string) []string
}

// ***
// Пакет для модульных тестов логики, зависящей от БД
type DB interface {
	GetUser(id int) (*User, error)
	AddUser(user *User) (*User, error)
}

type User struct {
	ID    int
	Email string
}

type MemDB struct {
	users  map[int]*User
	nextID int
}

func (m *MemDB) GetUser(id int) (*User, error) {
	return m.users[id], nil
}

func (m *MemDB) AddUser(user *User) (*User, error) {
	user.ID = m.nextID
	m.users[user.ID] = user
	m.nextID++
	return user, nil
}

// Интерфейсы стандартной библиотеки
func stdLib() {
	var a any

	switch val := a.(type) {
	// Обычно функциям следует принимать интерфейсы как параметры,
	// но возвращать конкретные типы.
	// Интерфейс `error` является исключением из правил.
	case error:
		log.Info().Msgf("var a is error: %v", val.Error())
	case fmt.Stringer:
		log.Info().Msg(val.String())
	case io.Reader:
		val.Read(nil)
	case io.Writer:
		val.Write(nil)
	case io.Closer:
		val.Close()
	case json.Marshaler:
		val.MarshalJSON()
	case json.Unmarshaler:
		val.UnmarshalJSON(nil)
	default:
		log.Info().Msg("var a is none of the known types")
	}
}
