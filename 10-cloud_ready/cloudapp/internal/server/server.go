package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"go-masters/10-cloud_ready/cloudapp/internal/config"
	"go-masters/10-cloud_ready/cloudapp/internal/db"
	"go-masters/10-cloud_ready/cloudapp/internal/db/postgres"
	"go-masters/10-cloud_ready/cloudapp/internal/metrics"
	"go-masters/10-cloud_ready/cloudapp/internal/models"
	"go-masters/10-cloud_ready/cloudapp/internal/telemetry"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	cfg    *config.Cfg
	router *chi.Mux
	server *http.Server
	db     db.DB
}

func New(cfg *config.Cfg) (*Server, error) {
	r := chi.NewRouter()

	db, err := postgres.New(cfg.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации БД: %w", err)
	}

	s := Server{
		cfg:    cfg,
		router: r,
		server: &http.Server{
			Addr:         fmt.Sprintf(":%v", cfg.Port),
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		db: db,
	}

	s.endpoints()

	return &s, nil
}

func (s *Server) endpoints() {
	// Настройка middleware
	s.router.Use(
		middleware.RequestID,                 // Добавляет X-Request-Id в заголовки
		telemetry.TracingMiddleware,          // OpenTelemetry трейсинг
		metrics.PrometheusMiddleware,         // Метрики Prometheus
		RequestLoggerMiddleware(&log.Logger), // Логирование запросов
		middleware.Recoverer,                 // Восстановление после паник
	)

	// Эндпоинты pprof
	// http://localhost:8080/debug/pprof/
	s.router.Get("/debug/pprof/", pprof.Index)
	s.router.Get("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.Get("/debug/pprof/profile", pprof.Profile)
	s.router.Get("/debug/pprof/symbol", pprof.Symbol)
	s.router.Get("/debug/pprof/trace", pprof.Trace)
	s.router.Get("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	s.router.Get("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
	s.router.Get("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.router.Get("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
	s.router.Get("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
	s.router.Get("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	// HealthCheck - статус системы
	s.router.Get("/health", healthHandler)

	// Эндпоинт для Prometheus
	s.router.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Инициализация маршрутов
	s.router.Post("/albums", s.addAlbumHandler)
	s.router.Get("/albums", s.listAlbumsHandler)
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

func (s *Server) addAlbumHandler(w http.ResponseWriter, r *http.Request) {
	span := trace.SpanFromContext(r.Context())
	defer span.End()

	log.Info().Msg("Обработка запроса addAlbum")
	span.AddEvent("Обработка запроса addAlbum")

	var req models.Album
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось декодировать запрос")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.db.AddAlbum(r.Context(), req)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось добавить альбом в БД")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) listAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx)
	defer span.End()

	log.Info().Msg("Обработка запроса listAlbums")
	span.AddEvent("Обработка запроса listAlbums")

	albums, err := s.db.ListAlbums(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "не удалось получить альбомы")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
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
