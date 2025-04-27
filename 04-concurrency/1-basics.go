package concurrency

import (
	"fmt"
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

// Запуск горутин.
func goroutine() {
	f := func(i int) {
		fmt.Printf("%v\n", i)
	}

	i := 10

	go f(i)

	go func() {
		f(i)
	}()

	go func() {
		fmt.Printf("%v\n", i)
	}()

	time.Sleep(1 * time.Second)
}

func mutex() {
	// Мьютекс используется для защиты общей памяти и
	// избежания "гонки".
	mu := &sync.Mutex{}
	var i int

	const gouroutinesNum = 100_000
	wg := &sync.WaitGroup{}

	unsafeIncrement := func() {
		i++
		wg.Done()
	}

	wg.Add(gouroutinesNum)
	for range gouroutinesNum {
		go unsafeIncrement()
	}

	wg.Wait()
	log.Info().Msgf("unsafeIncrement(): i = %v", i)

	i = 0

	safeIncrement := func() {
		mu.Lock()
		defer mu.Unlock()

		i++
		wg.Done()
	}

	wg.Add(gouroutinesNum)
	for range 10_0000 {
		go safeIncrement()
	}

	wg.Wait()
	log.Info().Msgf("safeIncrement: i = %v", i)
}

func channel() {
	// Каналы используются для обмена сообщениями.
	// Если есть канал - должна быть отдельная горутина
	// для чтения или записи.

	// Канал обязательно надо создать с помощью `make(...)``
	// Обычно используются небуферизованные каналы.
	ch := make(chan string)

	go func() {
		message := <-ch
		log.Info().Msgf("Дополнительный поток получил сообщение: %v. Начинаю работу", message)
		close(ch)
	}()

	ch <- "'Старт!'"

	<-ch
	log.Info().Msg("Основной поток получил сообщение от дополнительного")
}
