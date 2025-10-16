package bsgostuff_infrastructure

import (
	"context"
	"fmt"

	bsgostuff_config "github.com/beavernsticks/go-stuff/config"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Return new Postgresql db instance
func NewPostgresDB(config bsgostuff_config.PostgreSQL) (*pgxpool.Pool, error) {
	dataSourceName := config.ConnectionString

	if dataSourceName == "" {
		dataSourceName = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			config.Host,
			config.Port,
			config.User,
			config.DBName,
			config.Password,
		)
	}

	pool, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		return nil, err
	}

	pool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

// MustNewPostgresDB создает адаптер или паникует при ошибке
func MustNewPostgresDB(cfg bsgostuff_config.PostgreSQL) *pgxpool.Pool {
	pool, err := NewPostgresDB(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize PostgreSQL connection: %w", err))
	}
	return pool
}
