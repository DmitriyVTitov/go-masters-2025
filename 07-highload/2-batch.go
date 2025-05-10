package highload

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbURL        = "postgres://user:password@localhost:5432/dbname"
	insertQuery  = "INSERT INTO test_table (id, name, value, created_at) VALUES ($1, $2, $3, $4)"
	rowsToInsert = 100_000
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
			log.Printf("Failed to insert row %d: %v\n", i, err)
		}
	}

	elapsed := time.Since(start)
	log.Printf("Inserted %d rows one by one in %s\n", rowsToInsert, elapsed)
}

func batchInsert() {
	ctx := context.Background()

	start := time.Now()

	// Используем пакетную вставку
	batch := &pgx.Batch{}
	count := 0

	for i := 1; i <= rowsToInsert; i++ {
		batch.Queue(
			"INSERT INTO test_table (id, name, value, created_at) VALUES ($1, $2, $3, $4)",
			i, fmt.Sprintf("Item %d", i), float64(i)*1.1, time.Now(),
		)
		count++

		// Отправляем пакет при достижении batchSize или в конце
		if count%batchSize == 0 || i == rowsToInsert {
			br := pool.SendBatch(ctx, batch)
			_, err := br.Exec()
			if err != nil {
				log.Printf("Batch insert failed: %v\n", err)
			}
			err = br.Close()
			if err != nil {
				log.Printf("Batch close failed: %v\n", err)
			}
			batch = &pgx.Batch{} // Создаем новый пакет
		}
	}

	elapsed := time.Since(start)
	log.Printf("Inserted %d rows using batches in %s\n", rowsToInsert, elapsed)
}
