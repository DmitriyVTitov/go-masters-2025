package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go-masters/10-cloud_ready/cloudapp/internal/config"
	"go-masters/10-cloud_ready/cloudapp/internal/server"

	"github.com/rs/zerolog/log"
)

func main() {
	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем конфигурацию
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка при загрузке конфигурации")
	}

	// Инициализируем сервер
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка при инициализации сервера")
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(ctx); err != nil {
			log.Error().Err(err).Msg("Ошибка при запуске сервера")
			cancel()
		}
	}()

	// Ожидаем сигналы ОС для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Получен сигнал на завершение работы")
	cancel()

	// Ожидаем завершения работы сервера
	<-ctx.Done()
	log.Info().Msg("Сервер успешно остановлен")
}
