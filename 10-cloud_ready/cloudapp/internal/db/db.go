package db

import (
	"context"
	"go-masters/10-cloud_ready/cloudapp/internal/models"
)

type DB interface {
	AddAlbum(context.Context, models.Album) error
	ListAlbums(context.Context) ([]models.Album, error)
}
