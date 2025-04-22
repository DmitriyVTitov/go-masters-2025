package generics

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type myInt int

type Numbers interface {
	~int | int8 | int16 | int32 | int64 |
		float32 | float64
}

func addNumbers[T Numbers](a, b T) T {
	return a + b
}

func addIntegers[T constraints.Integer](a, b T) T {
	return a + b
}

type T1 struct{}

type T2 struct{}

func (i myInt) String() string {
	return fmt.Sprintf("%v", int(i))
}

type Mixed interface {
	T1 | T2 | ~int | ~struct{ bool string }

	fmt.Stringer
}

func useMixed[T Mixed](t T) {
	t.String()
}
