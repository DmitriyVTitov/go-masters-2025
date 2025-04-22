package generics

import (
	"encoding/json"
	"net/http"
	"sync"
)

// ***
// Универсальная функция для раскодирования тела HTTP-запроса
// из формата JSON в любой тип данных.
// ***

type Validator interface {
	Validate() error
}

func decodeAndValidateRequest[T Validator](r *http.Request) (*T, error) {
	var req T

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	err = req.Validate()
	if err != nil {
		return nil, err
	}

	return &req, nil
}

// ***
// Q - очередь FIFO для объектов произвольного типа данных.
// ***
type Q[T any] struct {
	mu       *sync.Mutex
	elements []*T
}

func New[T any]() *Q[T] {
	q := Q[T]{
		mu: &sync.Mutex{},
	}

	return &q
}

// Enqueue добавляет объект в конец очереди.
func (q *Q[T]) Enqueue(el *T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.elements = append(q.elements, el)
}

// Dequeue возвращает первый объект в очереди.
func (q *Q[T]) Dequeue() *T {
	q.mu.Lock()
	defer q.mu.Unlock()

	var item *T

	if len(q.elements) > 0 {
		item = q.elements[0]
		q.elements = q.elements[1:]
	}

	return item
}

// Len возвращает количество элементов в очереди.
func (q *Q[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.elements)
}

// ***
// FanIn объединяет несколько входных каналов в один выходной.
// ***
func FanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup

	// Определяем функцию, которая читает данные из канала и передаёт их в выходной канал.
	worker := func(c <-chan int) {
		defer wg.Done()
		for value := range c {
			out <- value
		}
	}

	// Для каждого входного канала запускаем горутину.
	wg.Add(len(channels))
	for _, ch := range channels {
		go worker(ch)
	}

	// Отдельная горутина закрывает выходной канал, когда все входные каналы завершены.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
