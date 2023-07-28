package db

import (
	"database/sql"
	"fmt"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/context"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func InitPGXPool(cfg *PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	return InitPGXPoolFromDSN(dsn)
}

func InitPGXPoolFromDSN(dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	max := runtime.NumCPU() * 4
	config.MaxConns = int32(max)

	return pgxpool.NewWithConfig(context.Background(), config)
}

func InitSQL(cfg *PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	_, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return sql.Open("pgx", dsn)
}
