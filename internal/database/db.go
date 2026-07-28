package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucho2027/blog/internal/database/sqlc"
)

type DB struct {
	pool *pgxpool.Pool
	q    *sqlc.Queries
}

type DBConfig struct {
	DatabaseURL string
	MaxConn     int
}

func New(ctx context.Context, dbConf DBConfig) (*DB, error) {
	if dbConf.DatabaseURL == "" {
		return nil, fmt.Errorf("DatabaseURL should not be empty")
	}
	if dbConf.MaxConn < 0 {
		return nil, fmt.Errorf("MaxConn should be >= 0")
	}
	if dbConf.MaxConn == 0 {
		dbConf.MaxConn = 10
	}
	poolConfig, err := pgxpool.ParseConfig(dbConf.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DatabaseURL: %w", err)
	}
	poolConfig.MaxConns = int32(dbConf.MaxConn)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to createPool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping pool %w", err)
	}
	queries := sqlc.New(pool)
	return &DB{
		pool: pool,
		q:    queries,
	}, nil
}

func (db *DB) Queries() *sqlc.Queries {
	return db.q
}

func (db *DB) Close() {
	if db == nil {
		return
	}
	if db.pool == nil {
		return
	}
	db.pool.Close()
	db.pool = nil
}

func (db *DB) WithinTx(ctx context.Context, txBound func(q *sqlc.Queries) error) error {
	if db == nil || db.pool == nil {
		return fmt.Errorf("db not initialized")
	}
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("unable to start transaction %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			tx.Rollback(ctx)
			fmt.Println("error rolling back") // TODO: We should probably handle rollback errors here.
		}
	}()

	txQueries := db.q.WithTx(tx)

	err = txBound(txQueries)
	if err != nil {
		return fmt.Errorf("transaction callback failed %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit failed %w", err)
	}
	committed = true
	return nil
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
