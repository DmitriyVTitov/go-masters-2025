package generics

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

// ***
// Допустим, мы хотим написать алгоритм, который суммирует
// все элементы слайса.
// Мы предполагаем, что в слайсе могут быть только числа.
//
// Без дженериков нам придётся писать несколько функций
// для каждого типа данных.
// ***

// sumSliceInts суммирует все элементы слайса int.
func sumSliceInts(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum
}

// sumSliceFloats суммирует все элементы слайса float64.
func sumSliceFloats(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum
}

// Либо использовать пустой интерфейс и приведение типов.
// Но в таком случае мы теряем безопасность типов.
// И нам придётся писать много кода для обработки
// различных типов данных.

// sumSliceAny суммирует все элементы слайса любого типа.
func sumSliceAny(slice any) float64 {
	var sum float64
	switch s := slice.(type) {
	case []int:
		for _, v := range s {
			sum += float64(v)
		}
	case []float64:
		for _, v := range s {
			sum += v
		}
	default:
		// Нелогично возвращать ошибку в таком алгоритме.
		// Но что делать!?
		panic("unsupported type")
	}

	// Мы всегда возвращаем float64, даже для int.
	return sum
}

// С дженериками мы можем написать одну функцию,
// которая будет работать с любым типом данных.
// При этом мы сохраняем безопасность типов.
//
// T - это параметр типа.
// any - это ограничение типа.
// В данном случае мы говорим, что T может быть любым типом.

// sumSlice суммирует все элементы слайса любого из указанных типов.
func sumSlice[T int | float64](s []T) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// ***
// Что быстрее: обобщенная функция, или использование `any`?
// ***
