package postgres

import (
	"context"
	"go-masters/10-cloud_ready/cloudapp/internal/models"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(connstr string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	pg := Postgres{pool: pool}

	err = pg.applyMigrations()
	if err != nil {
		return nil, err
	}

	return &pg, nil
}

func (pg *Postgres) applyMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	parent := filepath.Dir(wd)
	migrationsFS := os.DirFS(parent)
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pg.pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return db.Close()
}

func (pg *Postgres) AddAlbum(ctx context.Context, album models.Album) error {
	_, err := pg.pool.Exec(
		ctx,
		"INSERT INTO albums (id  artist, title, year) VALUES ($1, $2, $3, $4)",
		album.ID,
		album.Artist,
		album.Title,
		album.Year)
	return err
}

func (pg *Postgres) ListAlbums(ctx context.Context) ([]models.Album, error) {
	rows, err := pg.pool.Query(ctx, "SELECT id, artist, title, year FROM albums")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(
			&album.ID,
			&album.Artist,
			&album.Title,
			&album.Year,
		); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}
