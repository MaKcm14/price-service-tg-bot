package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func New(ctx context.Context, dsn string, log *slog.Logger) (PostgreSQLRepo, error) {
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Error(fmt.Sprintf("error of connecting to the DB: %v", err))
		return PostgreSQLRepo{}, ErrDBConnection
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		log.Error(fmt.Sprintf("error of connecting to the DB: %v", err))
		return PostgreSQLRepo{}, ErrDBConnection
	}

	return PostgreSQLRepo{
		pool:   pool,
		logger: log,
	}, nil
}

// IsUserExists checks that the user with the chatID exists at the DB.
func (p PostgreSQLRepo) IsUserExists(ctx context.Context, chatID int64) (bool, error) {
	const op = "postgres.check-user"

	rows, err := p.pool.Query(ctx, "SELECT id FROM users WHERE telegram_id=$1", chatID)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return false, ErrQueryExec
	}
	defer rows.Close()

	return rows.Next(), nil
}

// AddUser adds the user with the chatID to the DB.
func (p PostgreSQLRepo) AddUser(ctx context.Context, chatID int64) error {
	const op = "postgres.adding-user"

	_, err := p.pool.Exec(ctx, "INSERT INTO users (telegram_id) VALUES ($1)", chatID)

	if err != nil {
		p.logger.Warn(fmt.Sprintf("error of the %v: %v", op, err))
		return ErrQueryExec
	}

	return nil
}

func (p PostgreSQLRepo) Close() {
	p.pool.Close()
}
