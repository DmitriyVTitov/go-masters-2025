package testingdebug

import (
	"errors"
	"strings"
)

// ***
// Unit Tests
// ***

// Тестируемые функции: разворот строки с помощью двух
// указателей и с помощью strings.Builder (slice)

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func ReverseStringWithError(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty string")
	}
	return ReverseString(s), nil
}

// Для демонстрации правил для имен тестов.

func lowercaseName() {}

type T struct{}

func (t *T) Method() {}

// ***
// Benchmarks
// ***

// Если у нас есть альтернативный вариант реализации алгоритма,
// то можно померить и сравнить их производитеьность.
func ReverseStringSB(s string) string {
	var sb strings.Builder
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		sb.WriteRune(runes[i])
	}
	return sb.String()
}
