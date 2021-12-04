package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"greenlight.nickherrig.com/internal/data"
	"greenlight.nickherrig.com/internal/jsonlog"

	"github.com/jackc/pgx/v4/pgxpool"
)

// will be generated automatically at build time.
const version = "0.0.1"

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxConns    int
		minConns    int
		maxIdleTime string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config  config
	logger  *jsonlog.Logger
	storage data.Storage
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "port for api")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxConns, "db-max-conns", 50, "PostgreSQL max connections")
	flag.IntVar(&cfg.db.minConns, "db-min-conns", 25, "PostgreSQL min connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle connection time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "rate limiter requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "rate limiter bursts")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "rate limiter enabled")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config:  cfg,
		logger:  logger,
		storage: data.NewStorage(db),
	}

	// TODO Read about servermux godoc
	// https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/

	// TODO Read up on http.Server
	//https://pkg.go.dev/net/http
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"env":  cfg.env,
		"addr": srv.Addr,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connConf, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	connConf.MaxConns = int32(cfg.db.maxConns)
	connConf.MinConns = int32(cfg.db.minConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	connConf.MaxConnIdleTime = duration

	db, err := pgxpool.ConnectConfig(ctx, connConf)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}
