package highload

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Создаем пул соединений
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal().Msgf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	// Подготовка таблицы (опционально)
	_, err = pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS test_table (
		id SERIAL PRIMARY KEY,
		name TEXT,
		value FLOAT,
		created_at TIMESTAMP
	)`)
	if err != nil {
		log.Fatal().Msgf("Unable to create table: %v\n", err)
	}

	m.Run()
}

func Test_simpleInsert(t *testing.T) {
	simpleInsert()
}

func Test_batchInsert(t *testing.T) {
	batchInsert()
}
