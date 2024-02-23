package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/heetch/sqalx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Queryable = sqalx.Node

type DBGetter func(context.Context) Queryable

func Connect(ctx context.Context, username string, password string, host string, port string, database string) (func(context.Context) Queryable, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)

	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	node, err := sqalx.New(db, sqalx.SavePoint(true))
	if err != nil {
		return nil, fmt.Errorf("failed to create sqalx node: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := node.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close sqalx node")
		}

		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close db")
		}
	}()

	log.Debug().Msg("connected to database")
	return func(ctx context.Context) Queryable {
		if tx := TxFromContext(ctx); tx != nil {
			return tx
		}

		return node
	}, nil
}
