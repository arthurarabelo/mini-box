package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// User espelha uma linha da tabela users.
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Name         string
	Email        string
}

var (
	ErrUserNotFound = errors.New("usuário não encontrado")
	ErrUserExists   = errors.New("usuário já existe")
)

// GetUserByUsername busca um usuário pelo username.
// Usado pelo LoginHandler para validar credenciais.
func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var u User
	err := Pool.QueryRow(ctx,
		`SELECT id, username, password_hash, name, email FROM users WHERE username = $1`,
		username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Name, &u.Email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUser insere um novo usuário. passwordHash já deve vir hasheado
// (bcrypt) — esta função não faz hashing.
//
// Usada pelo RegisterHandler (exercício em auth/handler.go).
func CreateUser(ctx context.Context, username, passwordHash, name, email string) (*User, error) {
	var u User
	err := Pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, name, email)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, username, password_hash, name, email`,
		username, passwordHash, name, email,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Name, &u.Email)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return nil, ErrUserExists
		}
		return nil, err
	}
	return &u, nil
}
