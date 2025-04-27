package concurrency

import (
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
)

// ***
// Для канала можно указать ограничени: только чтение или только запись.
// ***

func readOnlyChan(ch <-chan int) {
	// ch <- 10 - запрещено!
	i := <-ch
	println(i)
}

func writeOnlyChan(ch chan<- int) {
	// i := <-ch - запрещено!
	ch <- 10
}

// Неблокирующая запись в канал.
func nonBlockingWrite() {
	ch := make(chan int)

	select {
	case ch <- 10:
		log.Info().Msg("Запись в канал прошла успешно")
	default:
		log.Info().Msg("Канал заполнен, запись невозможна")
	}

	close(ch)
}

// Шаблон "ожидание завершения горутин".
// Реализуется с помощью WaitGroup
func waitGroup() {
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			log.Info().Msgf("Горутина %v завершила работу", i)
		}(i)
	}

	wg.Wait()
	log.Info().Msg("Все горутины завершили работу")
}

// Шаблон "семафор".
// Шаблон можно реализовать с помощью буферизованного канала,
// где размер буфера определяет количество разрешений (permits).
// Это классический шаблон для ограничения количества одновременно
// выполняемых операций.
func sema() {
	sem := make(chan struct{}, 3) // Макс горутин
	var counter atomic.Int32

	var wg sync.WaitGroup
	wg.Add(10)

	for range 10 {
		sem <- struct{}{}
		go func() {
			defer func() {
				counter.Add(-1)
				<-sem
				wg.Done()
			}()
			counter.Add(1)
			log.Info().Msgf("сейчас запущено %v горутин", counter.Load())
		}()
	}

	wg.Wait()
}
