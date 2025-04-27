package concurrency

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Шаблон Worker Pool (Пул рабочих потоков) на Go, который
// позволяет эффективно распределять задачи между несколькими
// горутинами, ограничивая их количество для контроля нагрузки.

// Task — задача, которую будут выполнять воркеры.
type Task string

func workersPool() {
	// 1. Создаем каналы для задач и результатов.
	tasks := make(chan Task)
	results := make(chan string)

	// 2. Запускаем пул воркеров (например, 3 горутины).
	var wg sync.WaitGroup

	numWorkers := runtime.GOMAXPROCS(-1)
	for i := range numWorkers {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// 3. Выводим результаты.
	done := make(chan struct{})
	go func() {
		defer close(done)

		for result := range results {
			fmt.Println(result)
		}
	}()

	// 4. Отправляем задачи в канал.
	for i := 1; i <= 10; i++ {
		tasks <- Task(fmt.Sprintf("Задание №%d", i))
	}
	close(tasks) // Закрываем канал задач (воркеры завершатся после выполнения всех задач).

	// 5. Ждем завершения всех воркеров.
	wg.Wait()
	close(results) // Закрываем канал результатов.
	<-done         // Ждем завершения горутины для вывода результатов.
}

// worker — функция, обрабатывающая задачи.
func worker(id int, tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		// Имитируем обработку задачи (например, задержку).
		time.Sleep(1 * time.Second)
		results <- fmt.Sprintf("Воркер %d: обработал задачу %s", id, task)
	}
}

// ***
// Шблон Pipeline (Конвейер) на Go, где задачи проходят через цепочку этапов
// обработки, каждый из которых выполняется в отдельной горутине.
// Это особенно полезно для последовательной обработки данных с четким
// разделением ответственности между этапами.
// ***

// ***
// Обобщенная структура конвейера
// ***

// Stage — функция, обрабатывающая данные и передающая их в следующий этап.
type Stage[T any] func(<-chan T) <-chan T

// Pipeline — объединяет этапы в цепочку.
func Pipeline[T any](stages ...Stage[T]) Stage[T] {
	return func(input <-chan T) <-chan T {
		ch := input
		for _, stage := range stages {
			ch = stage(ch)
		}
		return ch
	}
}

func pipelineExample() {
	// 1. Генерация чисел (первый этап).
	gen := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 5; i++ {
				out <- i
				time.Sleep(100 * time.Millisecond) // Имитация работы.
			}
		}()
		return out
	}

	// 2. Определяем этапы обработки.
	multiply := func(factor int) Stage[int] {
		return func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n * factor
					time.Sleep(200 * time.Millisecond)
				}
			}()
			return out
		}
	}

	add := func(offset int) Stage[int] {
		return func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for n := range in {
					out <- n + offset
					time.Sleep(200 * time.Millisecond)
				}
			}()
			return out
		}
	}

	// 3. Собираем пайплайн: Генератор -> Умножение на 2 -> Добавление 10.
	pipeline := Pipeline(
		multiply(2),
		add(5),
	)

	// 4. Запускаем пайплайн.
	source := gen()
	result := pipeline(source)

	// 5. Выводим результаты.
	for num := range result {
		fmt.Printf("Результат: %d\n", num)
	}
}

// ***
// Пример использования пакета context в Go для управления временем
// выполнения горутины.
// ***

// Функция, которая выполняет работу и учитывает контекст
func doWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Операция отменена:", ctx.Err())
			return
		default:
			// Чтобы учитывать контекст при длительных операциях,
			// в горутине должен быть цикл и операция должна
			// состоять из небольших квантов работы.
			fmt.Println("Работа выполняется...")
			time.Sleep(500 * time.Millisecond) // Имитация работы
		}
	}
}

func useContext() {
	// Создаем контекст с тайм-аутом 2 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Запускаем горутину с управляемым контекстом
	go doWork(ctx)

	// Ожидаем завершения
	time.Sleep(3 * time.Second)
	fmt.Println("Основная функция завершена")
}
