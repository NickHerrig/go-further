package data

import (
	"time"
)

// structs meant to represent json should be exported
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // - directive ommits field from json
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   int32     `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
