package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreDB struct {
	Pool     *pgxpool.Pool
	DBConfig *pgxpool.Config
}

type Config struct {
	Dsn          string `env:"POSTGRES_DSN,required"`
	MaxOpenConns int32  `env:"POSTGRES_MAX_OPEN_CONN" envDefault:"25"`
	MaxIdleConns int    `env:"POSTGRES_MAX_IDLE_CONN" envDefault:"25"`
	MaxIdleTime  string `env:"POSTGRES_MAX_IDLE_TIME" envDefault:"15m"`
}

func New(ctx context.Context, config Config) (*PostgreDB, error) {
	dbConfig, err := pgxpool.ParseConfig(config.Dsn)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = config.MaxOpenConns

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConnIdleTime = duration

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgreDB{
		Pool:     pool,
		DBConfig: dbConfig,
	}, nil
}
