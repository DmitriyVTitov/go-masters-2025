package generics

import "time"

// ***
// Обобщенная структура.
// ***

// Результат измерения датчика.
type Measure[T Numbers] struct {
	// Время измерения.
	Timestamp time.Time
	// Название измеряемого показателя.
	Metric string
	// Значение измерения.
	Value T
}

func genericStruct() {
	var temperatureMeasure Measure[int]
	var humidityMeasure Measure[float64]

	temperatureMeasure.Metric = "temperature"
	temperatureMeasure.Value = 25
	temperatureMeasure.Timestamp = time.Now()

	humidityMeasure.Metric = "humidity"
	humidityMeasure.Value = 60.5
	humidityMeasure.Timestamp = time.Now()

	_, _ = temperatureMeasure, humidityMeasure
}

// ***
// Обобщенный интерфейс.
// ***

type GenericInterface[T any] interface {
	Print(T)
}

type GenericImplementation[T any] struct{}

func (*GenericImplementation[T]) Print(v T) {
	println(v)
}

// Обобщение переходит вверх по пути наследования.
// Если вложенная структура обобщенная, то и внешняя должна быть обобщенной.
type Inner[T any] struct {
	Value T
}

type Outer[T any] struct {
	inner Inner[T]
}

// Либо нужно определить конкретный тип.
type OuterConcrete struct {
	innerInt Inner[int]
}
