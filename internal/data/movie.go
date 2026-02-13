package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"skuknuraknu/internal/models"
)

// MovieModel defines the interface for movie-related database operations.
type MovieModel interface {
	Insert(movie *models.Movie) error
	Get(id int64) (*models.Movie, error)
	// Update(movie *models.Movie) error
	// Delete(id int64) error
}

// MySQLMovieModel wraps a sql.DB connection pool.
type MySQLMovieModel struct {
	DB *sql.DB
}

// Insert adds a new movie to the MySQL database.
func (m *MySQLMovieModel) Insert(movie *models.Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES (?, ?, ?, ?)
	`

	// Convert genres slice to JSON for storage
	genresJSON, err := json.Marshal(movie.Genres)
	if err != nil {
		return err
	}

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query.
	result, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Year, movie.Runtime, genresJSON)
	if err != nil {
		return err
	}

	// Retrieve the ID of the newly inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	movie.ID = id
	// In a real scenario, we might want to fetch the created_at timestamp from the DB,
	// but for simplicity and since MySQL doesn't have RETURNING clause like PostgreSQL,
	// we'll set it here or fetch it separately if needed.
	movie.Version = 1
	return nil
}

// Get retrieves a movie by its ID from the MySQL database.
func (m *MySQLMovieModel) Get(id int64) (*models.Movie, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = ?
	`

	var movie models.Movie
	var genres []byte // To handle JSON data from MySQL

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query and scan the result into the movie struct.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&genres,
		&movie.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}

	// Unmarshal the JSON genres data into the slice.
	if err := json.Unmarshal(genres, &movie.Genres); err != nil {
		return nil, err
	}

	return &movie, nil
}
