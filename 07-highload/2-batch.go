package highload

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const (
	dbURL       = "postgres://postgres:password@localhost:5432/db"
	insertQuery = `
		INSERT INTO
			test_table (id, name, value, created_at)
		VALUES
			($1, $2, $3, $4);`
	rowsToInsert = 10_000
	batchSize    = 1000 // Размер пакета
)

var (
	pool *pgxpool.Pool
)

func simpleInsert() {
	ctx := context.Background()

	start := time.Now()

	// Вставка строк по одной
	for i := 1; i <= rowsToInsert; i++ {
		_, err := pool.Exec(ctx, insertQuery, i, fmt.Sprintf("Item %d", i), float64(i)*1.1, time.Now())
		if err != nil {
			log.Fatal().Msgf("Failed to insert row %d: %v\n", i, err)
		}
	}

	elapsed := time.Since(start)
	log.Info().Msgf("Добавлено %d строк по одной за %s\n", rowsToInsert, elapsed)
}

func batchInsert() {
	ctx := context.Background()

	start := time.Now()

	// Используем пакетную вставку
	batch := &pgx.Batch{}
	count := 0

	for i := 1; i <= rowsToInsert; i++ {
		batch.Queue(
			insertQuery, i, fmt.Sprintf("Item %d", i), float64(i)*1.1, time.Now(),
		)
		count++

		// Отправляем пакет при достижении batchSize или в конце
		if count%batchSize == 0 || i == rowsToInsert {
			br := pool.SendBatch(ctx, batch)
			_, err := br.Exec()
			if err != nil {
				log.Fatal().Msgf("Batch insert failed: %v\n", err)
			}
			err = br.Close()
			if err != nil {
				log.Fatal().Msgf("Batch close failed: %v\n", err)
			}
			batch = &pgx.Batch{} // Создаем новый пакет
		}
	}

	elapsed := time.Since(start)
	log.Info().Msgf("Добавлено %d строк пакетами по %d за %s\n",
		rowsToInsert, batchSize, elapsed)
}
