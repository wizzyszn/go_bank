package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

/*
*
What is a connection pool?

A connection pool is a cache of database connections that can be reused. Instead of creating a new connection every time you need one, the pool maintains a set of ready-to-use connections.

What does it do?

Reuses connections - when you finish using a connection, it goes back to the pool rather than being closed
Manages the lifecycle - automatically creates connections up to a limit and closes idle ones
Handles concurrency - multiple goroutines can safely request connections simultaneously
Maintains state per connection - if your database has per-connection state (sessions, transactions), it stays isolated
Why does it exist?

Creating a new database connection is expensive:

TCP handshake overhead
Authentication/login process
Memory allocation
Network latency
Without a pool, every query would pay these costs. With a pool:

Performance - reusing connections is much faster than creating new ones
Resource efficiency - limits total connections to your database (prevents overwhelming the server)
Better scalability - allows many concurrent requests to share a limited number of connections
In Go's sql.DB:

The pool automatically manages connections for you
You can tune it with SetMaxIdleConns() and SetMaxOpenConns()
It's safe for concurrent use - multiple goroutines can query simultaneously
When you call Begin() for a transaction, it gets one dedicated connection from the pool until you commit/rollback
This is why Go provides *sql.DB as a singleton—you create one pool at startup and reuse it for the lifetime of your application.
*/
type DB struct {
	*sql.DB
}

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func New(cfg Config) (*DB, error) {

	db, err := sql.Open("postgres", cfg.DSN)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)

	}

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting databse %w", err)
	}

	log.Println("Databse connection established")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	log.Println("Closing database connetion")
	return db.DB.Close()
}

func (db *DB) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed")
	}
	return nil
}

func (db *DB) Stats() sql.DBStats {
	return db.DB.Stats()
}
func NewConfig(dsn string) Config {
	return Config{
		DSN:             dsn,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}
