package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

// RepositoryPostgres type replies for accessing to postgres database
type RepositoryPostgres struct {
	DB *pgxpool.Pool
}
