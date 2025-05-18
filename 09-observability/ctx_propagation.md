# Пример распространения контекста с использованием OpenTelemetry (OTEL) в Go

Вот полный пример, демонстрирующий распространение контекста с использованием OpenTelemetry в Go. В примере показано, как создавать spans, передавать контекст между сервисами и извлекать контекст из входящих запросов.

## Пример: HTTP-сервис с распространением контекста

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

...

// handleRequest - HTTP-обработчик, который создает span и вызывает другой сервис
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Извлекаем контекст из входящего запроса
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

	// Начинаем новый span
	ctx, span := tracer.Start(ctx, "handleRequest")
	defer span.End()

	// Добавляем атрибуты к span
	span.SetAttributes(
		semconv.HTTPMethod(r.Method),
		semconv.HTTPURL(r.URL.String()),
	)

	// Имитируем работу
	time.Sleep(100 * time.Millisecond)

	// Вызываем другой сервис (может быть HTTP, gRPC и т.д.)
	callAnotherService(ctx)

	w.Write([]byte("Запрос успешно обработан"))
}

// callAnotherService демонстрирует передачу контекста в другой сервис
func callAnotherService(ctx context.Context) {
	// Начинаем новый span как дочерний для входящего контекста
	_, span := tracer.Start(ctx, "callAnotherService")
	defer span.End()

	// Имитируем вызов другого сервиса
	time.Sleep(50 * time.Millisecond)

	// Добавляем событие в span
	span.AddEvent("Вызван другой сервис")

	// Можно передать контекст дальше в другие сервисы
	// Для HTTP создаем заголовки с распространением:
	headers := make(http.Header)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))

	// Теперь 'headers' содержит traceparent и другие заголовки распространения
	// Эти заголовки используются при выполнении HTTP-запросов к другим сервисам
}

// makeRequest демонстрирует распространение контекста при выполнении исходящих запросов
func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	// Начинаем span для этой операции
	ctx, span := tracer.Start(ctx, "makeRequest")
	defer span.End()

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Внедряем контекст трассировки в заголовки запроса
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Выполняем запрос
	client := http.Client{}
	return client.Do(req)
}
```

## Основные концепции в примере:
   - Извлечение контекста из входящих HTTP-запросов
   - Внедрение контекста в исходящие HTTP-запросы

Этот пример показывает как извлечение контекста на стороне сервера, так и внедрение контекста на стороне клиента, демонстрируя полный сценарий распространения. Трассировки будут отображаться в интерфейсе Jaeger (обычно по адресу http://localhost:16686).