package errs

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func logError() {

	// Часто бывает нужно в конце какой-то функции выполнить
	// в отдельном потоке некоторые завершающие действия:
	// - отправить уведомление о результате операции
	// - закрыть какие-то ресурсы и т.д.
	// При этом результат таких завершающих операций не может считаться
	// результатом работы самой функции.
	// В таком случае ошибку можно только залогировать.

	ch := make(chan struct{})

	go func() {
		defer close(ch)

		err := sendMessage()
		if err != nil {
			log.Error().Err(err).Msg("не повезло")
		} else {
			log.Info().Msg("повезло")
		}
	}()

	<-ch
}

func sendMessage() error {
	if rand.IntN(10) < 5 {
		return errors.New("неудачное число")
	}

	return nil
}

func errChan() {
	errCh := make(chan error)

	wg := sync.WaitGroup{}
	wg.Add(100)
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for i := range 100 {
		go func() {
			goroutine(i, errCh)
			wg.Done()
		}()
	}

	for err := range errCh {
		log.Error().Err(err).Msg("ошибка из канала")
	}
}

func goroutine(index int, errCh chan error) {
	if rand.N(100) > 90 {
		errCh <- fmt.Errorf("поток №%v: не повезло", index)
	}
}

func errGroup() {
	eg, ctx := errgroup.WithContext(context.Background())

	for range 10 {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				log.Error().Err(ctx.Err()).Msg("выполнение горутины прервано")
				return ctx.Err()
			default:
				if rand.N(100) > 80 {
					return errors.New("не повезло")
				} else {
					log.Info().Msg("горутина выполнилась успешно")
					return nil
				}
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Error().Err(err).Msg("ошибка из errgroup")
	}
}
