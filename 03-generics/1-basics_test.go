package generics

import (
	"testing"

	"github.com/rs/zerolog/log"
)

func Test_sumSliceXxx(t *testing.T) {
	ints := []int{1, 2, 3}
	floats := []float64{1, 2, 3}

	log.Info().Msgf("sumSliceInts() = %v", sumSliceInts(ints))
	log.Info().Msgf("sumSliceFloats() = %v", sumSliceFloats(floats))
	log.Info().Msgf("sumSliceAny(ints) = %v", sumSliceAny(ints))
	log.Info().Msgf("sumSliceAny(floats) = %v", sumSliceAny(floats))
	log.Info().Msgf("sumSlice(ints) = %v", sumSlice(ints))
	log.Info().Msgf("sumSlice(floats) = %v", sumSlice(floats))
}
