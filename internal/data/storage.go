package data

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("Woops, we coudn't find that record")
	ErrEditConflict   = errors.New("Woops, there are an edit conlfict")
)

type Storage struct {
	Movies MovieStorage
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Movies: MovieStorage{DB: db},
	}
}
