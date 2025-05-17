package highload

import (
	"testing"
)

// go test -gcflags="-m" -run XXX

func Test_woMake(t *testing.T) {
	woMake()
}

func Test_withMake(t *testing.T) {
	withMake()
}

func Benchmark_woMake(b *testing.B) {
	for b.Loop() {
		woMake()
	}
}

func Benchmark_withMake(b *testing.B) {
	for b.Loop() {
		withMake()
	}
}

func Test_usePool(t *testing.T) {
	usePool()
}

func Benchmark_existsInslice(b *testing.B) {
	slice := make([]int, 1000)
	for i := range 1000 {
		slice[i] = i
	}

	for b.Loop() {
		existsInslice(slice, 500)
	}
}

func Benchmark_existsInMap(b *testing.B) {
	m := make(map[int]struct{}, 1000)
	for i := range 1000 {
		m[i] = struct{}{}
	}

	for b.Loop() {
		existsInMap(m, 500)
	}
}

func BenchmarkBruteForce(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	target := 1999 // Максимально возможная сумма

	b.ResetTimer()
	for b.Loop() {
		twoSumBruteForce(nums, target)
	}
}

func BenchmarkHashMethod(b *testing.B) {
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i
	}
	target := 1999

	b.ResetTimer()
	for b.Loop() {
		twoSumHash(nums, target)
	}
}
