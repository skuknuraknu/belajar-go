package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"skuknuraknu/internal/config"
	"skuknuraknu/internal/data"

	_ "github.com/go-sql-driver/mysql"
)

const version = "1.0.0"

// application holds the dependencies for the HTTP handlers.
type application struct {
	config     *config.Config
	logger     *log.Logger
	movieModel data.MovieModel
}

func main() {
	// Load application configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize a logger.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Open a database connection pool.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Printf("database connection pool established")

	// Initialize the movie model with the DB connection.
	// We use the MySQL implementation now.
	movieModel := &data.MySQLMovieModel{DB: db}

	// Initialize an application struct with dependencies.
	app := &application{
		config:     cfg,
		logger:     logger,
		movieModel: movieModel,
	}

	// Initialize an HTTP router.
	router := app.routes()

	// Configure the HTTP server.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the server.
	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// openDB opens a database connection with the given configuration.
// It pings the database to verify the connection is established.
func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		return nil, err
	}

	// Set a timeout for the initial connection verification.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping the database to ensure connection is valid.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
