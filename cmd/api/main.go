// CMPS3162 Lab 3 demonstrating writeJSON and readJSON.
package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/andreshungbz/lab4-database-crud/internal/data"
	"github.com/andreshungbz/lab4-database-crud/internal/vcs"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

// config stores the API server configuration.
type config struct {
	port int    // API server port
	env  string // (development|staging|production)
	db   struct {
		dsn string // data source name
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// application holds the dependencies for the HTTP handlers, helpers, middleware,
// etc. so that they are all accessible through dependency injection.
type application struct {
	config config
	logger *slog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	var cfg config
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// FLAGS

	// server flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// database flags
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	// rate-limiter flags
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// version flag
	displayVersion := flag.Bool("version", false, "Display program version")

	flag.Parse()

	// display program version and exit if the version flag was passed
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	// DATABASE

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Database connection pool established")

	// METRICS

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	// APPLICATION

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// start the API server
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB connects to the PostgreSQL database using the provided DSN and
// and returns a pointer to a handler to that database.
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// test the connection with a ping
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
