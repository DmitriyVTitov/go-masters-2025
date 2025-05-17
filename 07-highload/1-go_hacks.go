package highload

import (
	"os"
	"sync"
	"time"

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
// Миниимизация аллокаций памяти
// ***

// Плохо: аллокация при каждом append
func woMake() {
	var slice []int
	for i := range 1000 {
		slice = append(slice, i) // Может вызывать копирование
	}
}

// Хорошо: предварительное выделение
func withMake() {
	slice := make([]int, 0, 1000) // capacity = 1000
	for i := range 1000 {
		slice = append(slice, i) // Без реаллокаций
	}

	// То же самое для хэш-таблиц
	//
	// m := make(map[string]int, 100) // Задаём начальный размер.
	//_ = m
}

// ***
// Использование конечных пулов для переиспользования объектов
// ***

func usePool() {
	var mu sync.Mutex
	var maxPoolSize int

	// Создаем пул с дорогостоящими объектами
	var pool = sync.Pool{
		// Конструктор объекта
		New: func() any {
			// Обновляем максимальный размер пула
			mu.Lock()
			maxPoolSize++
			mu.Unlock()

			// Создаем переиспользуемый ресурс (буффер)
			return make([]byte, 1024*1024)
		},
	}

	const N = 100
	var wg sync.WaitGroup
	wg.Add(N)
	for i := range N {
		go func() {
			defer wg.Done()

			// Получаем объект из пула (или создаем новый)
			buf := pool.Get()
			if b, ok := buf.([]byte); ok {
				log.Info().Msgf("поток %v: получен буффер длиной %v", i, len(b))
			}

			// Возвращаем объект в пул
			pool.Put(buf)
		}()

		// Можно поменять тут задержку, размер пула будет также меняться
		time.Sleep(time.Millisecond * 3)
	}
	wg.Wait()

	log.Info().Msgf("максимальный размер пула: %v", maxPoolSize)
}

// Подбор оптимальных структур данных под задачу
func existsInslice[T comparable](slice []T, value T) bool {
	for i := range len(slice) {
		if slice[i] == value {
			return true
		}
	}

	return false
}

func existsInMap[T comparable](m map[T]struct{}, value T) bool {
	_, ok := m[value]
	return ok
}

// ***
// Снижение оценки сложности алгоритма
// ***

// Задача: Найти все пары чисел в массиве, которые в сумме дают target

// O(n²) - перебор всех возможных пар
func twoSumBruteForce(nums []int, target int) [][2]int {
	var result [][2]int
	n := len(nums)

	for i := range n {
		for j := i + 1; j < n; j++ {
			if nums[i]+nums[j] == target {
				result = append(result, [2]int{nums[i], nums[j]})
			}
		}
	}
	return result
}

// O(n) - с использованием хэш-таблицы
func twoSumHash(nums []int, target int) [][2]int {
	var result [][2]int
	seen := make(map[int]bool)

	for _, num := range nums {
		complement := target - num
		if seen[complement] {
			result = append(result, [2]int{complement, num})
		}
		seen[num] = true
	}
	return result
}
