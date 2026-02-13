package models

import (
	"time"
)

// Movie represents a single movie record.
// It includes fields for an ID, title, release year, runtime, and genres.
type Movie struct {
	ID        int64     `json:"id"`      // Unique identifier for the movie.
	CreatedAt time.Time `json:"-"`       // Timestamp when the record was created (not exposed in JSON).
	Title     string    `json:"title"`   // The title of the movie.
	Year      int32     `json:"year"`    // The release year of the movie.
	Runtime   int32     `json:"runtime"` // The runtime of the movie in minutes.
	Genres    []string  `json:"genres"`  // A slice of genres associated with the movie.
	Version   int32     `json:"version"` // The version number, used for optimistic concurrency control.
}
