package server

import (
	"context"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"time"

	"go-masters/09-observability/example/internal/metrics"
	"go-masters/09-observability/example/internal/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	router *chi.Mux
	server *http.Server
}

func New() *Server {
	r := chi.NewRouter()

	// Настройка middleware
	r.Use(
		middleware.RequestID,                 // Добавляет X-Request-Id в заголовки
		telemetry.TracingMiddleware,          // OpenTelemetry трейсинг
		metrics.PrometheusMiddleware,         // Метрики Prometheus
		RequestLoggerMiddleware(&log.Logger), // Логирование запросов
		middleware.Recoverer,                 // Восстановление после паник
	)

	// Эндпоинты pprof
	// http://localhost:8080/debug/pprof/
	r.Get("/debug/pprof/", pprof.Index)
	r.Get("/debug/pprof/cmdline", pprof.Cmdline)
	r.Get("/debug/pprof/profile", pprof.Profile)
	r.Get("/debug/pprof/symbol", pprof.Symbol)
	r.Get("/debug/pprof/trace", pprof.Trace)
	r.Get("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	r.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	r.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	r.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	r.Get("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	r.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	// HealthCheck - статус системы
	r.Get("/health", healthHandler)

	// Эндпоинт для Prometheus
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Инициализация маршрутов
	r.Get("/hello", helloHandler)

	return &Server{
		router: r,
		server: &http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	log.Info().Msg("Инициализация телеметрии")

	shutdown, err := telemetry.SetupOTelSDK(ctx, "http://localhost:4318")
	if err != nil {
		return err
	}
	defer func() {
		err = shutdown(ctx)
		if err != nil {
			log.Err(err).Msg("cannot shutdown OTel")
		}
	}()

	log.Info().Str("addr", s.server.Addr).Msg("Запуск HTTP сервера")

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		log.Info().Msg("Остановка HTTP сервера")
		if err := s.server.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Ошибка при остановке сервера")
		}
	}()

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Обработчики запросов

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	log.Info().Msg("Обработка запроса hello")
	span.AddEvent("Обработка запроса hello")

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, " + name + "!"))

	if rand.Intn(10) < 5 {
		span.SetStatus(codes.Error, "random error")
	}
}

// RequestLoggerMiddleware - middleware для логирования запросов
func RequestLoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)

			logger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration", time.Since(start)).
				Str("request_id", middleware.GetReqID(r.Context())).
				Msg("Обработан HTTP запрос")
		})
	}
}
