package pkg

import (
	"net/http"
)

// Типичный пакет на Go в ООП-стиле.
// В пакете хранится некторое состояние (attempts, httpclient),
// поэтому он оформляется в виде структуры (Service)
// с набором методов.

// Service - основная сущность.
type Service struct {
	// Количество попыток выполнения запроса.
	attempts int
	// Клиент для выполнения запросов.
	httpclient *http.Client
}

// New - базовый конструктор.
func New() *Service {
	s := Service{
		attempts:   3,
		httpclient: http.DefaultClient,
	}

	return &s
}

// Option - функциональная опция.
type Option func(*Service)

func WithAttempts(attempts int) Option {
	return func(s *Service) {
		s.attempts = attempts
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(s *Service) {
		s.httpclient = client
	}
}

// NewWithOptions - расширенный конструктор с дополнительными опциями.
func NewWithOptions(opts ...Option) *Service {
	s := New()
	for _, opt := range opts {
		opt(s)
	}

	return s
}
