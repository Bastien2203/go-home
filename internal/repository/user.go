package repository

import (
	"database/sql"
	"fmt"

	"github.com/Bastien2203/go-home/internal/core"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) (*UserRepository, error) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		password_hash TEXT,
		created_at DATETIME default current_timestamp
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Save(user *core.User) error {
	query := `
	INSERT OR REPLACE INTO users 
	(id, email, password_hash)
	VALUES (?, ?, ?)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*core.User, error) {
	query := `SELECT id, email, created_at, password_hash FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)
	return r.scanUser(row)
}

func (r *UserRepository) FindByID(id string) (*core.User, error) {
	query := `SELECT id, email, created_at, password_hash FROM users WHERE id = ?`

	row := r.db.QueryRow(query, id)
	return r.scanUser(row)
}

func (r *UserRepository) Count() (*int, error) {
	query := `SELECT COUNT(id) FROM users`

	row := r.db.QueryRow(query)

	var count int
	if err := row.Scan(&count); err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *UserRepository) scanUser(row Scanner) (*core.User, error) {
	var u core.User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.CreatedAt,
		&u.PasswordHash,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
