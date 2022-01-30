package data

import (
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("Woops, we coudn't find that record")
	ErrEditConflict   = errors.New("Woops, there was an edit conlfict")
)

type Storage struct {
	Movies MovieStorage
	Users  UserStorage
	Tokens TokenStorage
}

func NewStorage(db *pgxpool.Pool) Storage {
	return Storage{
		Movies: MovieStorage{DB: db},
		Users:  UserStorage{DB: db},
		Tokens: TokenStorage{DB: db},
	}
}
