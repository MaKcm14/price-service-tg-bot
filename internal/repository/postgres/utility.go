package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// postgresConfig defines the main postgres DB configuration.
type postgresConfig struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// user defines the main user's ids:
// chatID - user's tgID;
// UserID - user's DB id.
type userID struct {
	ChatID int64
	UserID int64
}

func newPostgresConfig(pool *pgxpool.Pool, log *slog.Logger) postgresConfig {
	return postgresConfig{
		pool:   pool,
		logger: log,
	}
}
