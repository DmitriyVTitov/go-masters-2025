package highload

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Создаем пул соединений
	p, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal().Msgf("не удалось создать пул соединений: %v\n", err)
	}
	defer p.Close()

	pool = p

	// Подготовка таблицы (опционально)
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS test_table (
		id SERIAL PRIMARY KEY,
		name TEXT,
		value FLOAT,
		created_at TIMESTAMP
	)`)
	if err != nil {
		log.Fatal().Msgf("не удалось создать таблицу: %v\n", err)
	}

	_, err = pool.Exec(ctx, `TRUNCATE TABLE test_table;`)
	if err != nil {
		log.Fatal().Msgf("не удалось очистить таблицу: %v\n", err)
	}

	m.Run()
}

func Test_simpleInsert(t *testing.T) {
	simpleInsert()
}

func Test_batchInsert(t *testing.T) {
	batchInsert()
}
