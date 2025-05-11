package main

import (
	"encoding/json"
	"io"
)

// Принцип подстановки Барбары Лисков (Liskov Substitution Principle - LSP)
// Если функция принимает на вход базовый тип, то должна принимать и
// производные типы.

// Интерфейс с контрактом
type Serializer interface {
	Serialize() []byte
}

type String string

func (s String) Serialize() []byte {
	return []byte(s)
}

type Album struct {
	Title String
	Year  uint
}

func (a Album) Serialize() []byte {
	b, err := json.Marshal(a)
	if err != nil {
		return nil
	}
	return b
}

// WriteObject сохраняет сериализованное представление объектов.
func WriteObject(w io.Writer, objects ...Serializer) error {
	for _, obj := range objects {
		b := obj.Serialize()
		_, err := w.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}
