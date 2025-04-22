package generics

import "time"

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

type GenericInterface struct {
}
