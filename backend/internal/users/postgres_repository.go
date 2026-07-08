package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Migrate() error {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  stats JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	return err
}

func (r *PostgresRepository) Create(user User) error {
	stats, err := json.Marshal(user.Stats)
	if err != nil {
		return err
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}
	_, err = r.db.Exec(
		`INSERT INTO users (id, username, email, password_hash, stats, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		stats,
		user.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *PostgresRepository) FindByEmail(email string) (User, error) {
	return r.findOne(`SELECT id, username, email, password_hash, stats, created_at FROM users WHERE email = $1`, email)
}

func (r *PostgresRepository) FindByID(id string) (User, error) {
	return r.findOne(`SELECT id, username, email, password_hash, stats, created_at FROM users WHERE id = $1`, id)
}

func (r *PostgresRepository) Update(user User) error {
	stats, err := json.Marshal(user.Stats)
	if err != nil {
		return err
	}
	result, err := r.db.Exec(
		`UPDATE users SET username = $1, email = $2, password_hash = $3, stats = $4 WHERE id = $5`,
		user.Username,
		user.Email,
		user.PasswordHash,
		stats,
		user.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresRepository) List() []User {
	rows, err := r.db.Query(`SELECT id, username, email, password_hash, stats, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return []User{}
	}
	defer rows.Close()

	result := []User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err == nil {
			result = append(result, user)
		}
	}
	return result
}

func (r *PostgresRepository) findOne(query string, args ...any) (User, error) {
	user, err := scanUser(r.db.QueryRow(query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (User, error) {
	var user User
	var statsRaw []byte
	if err := scanner.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &statsRaw, &user.CreatedAt); err != nil {
		return User{}, err
	}
	if err := json.Unmarshal(statsRaw, &user.Stats); err != nil {
		return User{}, err
	}
	if user.Stats.SolvedByCategory == nil {
		user.Stats.SolvedByCategory = map[string]int{"easy": 0, "medium": 0, "hard": 0, "nightmare": 0, "AI": 0}
	}
	return user, nil
}
