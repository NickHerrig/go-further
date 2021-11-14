package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"greenlight.nickherrig.com/internal/validator"
)

// structs meant to represent json should be exported
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // - directive ommits field from json
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` //using MarshalJSON method under-the-hood
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MovieStorage struct {
	DB *pgxpool.Pool
}

func (m MovieStorage) Insert(movie *Movie) error {
	query := `
		INSERT INTO MOVIES (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, movie.Genres}

	return m.DB.QueryRow(context.Background(), query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

}

func (m MovieStorage) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1`

	var movie Movie

	err := m.DB.QueryRow(context.Background(), query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&movie.Genres,
		&movie.Version,
	)
	// more information here about specific errors
	//https://github.com/jackc/pgx/issues/474#issuecomment-657538224
	if err != nil {
		switch {
		case err.Error() == "no rows in result set":
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil

}

func (m MovieStorage) Update(movie *Movie) error {

	query := `
		UPDATE movies
        	SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 
		WHERE id = $5
        	RETURNING version`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		movie.Genres,
		movie.ID,
	}

	return m.DB.QueryRow(context.Background(), query, args...).Scan(&movie.Version)

}

func (m MovieStorage) Delete(id int64) error {

	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1`

	result, err := m.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}
