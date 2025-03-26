package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// user defines the main user's ids:
// chatID - user's tgID;
// userID - user's DB id.
type userID struct {
	chatID int64
	userID int64
}

// postgresConfig defines the main postgres DB configuration.
type postgresConfig struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func newPostgresConfig(pool *pgxpool.Pool, log *slog.Logger) postgresConfig {
	return postgresConfig{
		pool:   pool,
		logger: log,
	}
}
