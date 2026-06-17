package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Connect abre o pool de conexões com o Postgres.
// DATABASE_URL pode ser sobrescrita por variável de ambiente; o valor
// default aqui bate com o docker-compose.yml do projeto.
func Connect(ctx context.Context) error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://minibox:minibox@localhost:5432/minibox?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return fmt.Errorf("abrindo pool de conexões: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("conectando ao postgres: %w", err)
	}

	Pool = pool
	return nil
}
