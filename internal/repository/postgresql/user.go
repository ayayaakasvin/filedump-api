package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"up-down-server/internal/lib/bcrypthashing"
)

func (p *PostgreSQL) RegisterUser(ctx context.Context, username, hashedPassword string) error {
	if exists, err := p.UsernameExists(username); err != nil {
		return err
	} else if exists {
		return errors.New("user already exists")
	}

	stmt, err := p.conn.PrepareContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, username, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) AuthentificateUser(ctx context.Context, username, password string) (int, error) {
	stmt, err := p.conn.PrepareContext(ctx, "SELECT user_id, password FROM users WHERE username = $1 LIMIT 1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var (
		userId         int
		hashedPassword string
	)

	if err = stmt.QueryRowContext(ctx, username).Scan(&userId, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(NotFound)
		}
		return 0, err
	}

	if err := bcrypthashing.ComparePasswordAndHash(password, hashedPassword); err != nil {
		return 0, errors.New(UnAuthorized)
	}

	return userId, nil
}

// UsernameExists checks if a record with the given username exists in the PostgreSQL database
func (p *PostgreSQL) UsernameExists(name string) (bool, error) {
	stmt, err := p.conn.Prepare("SELECT username FROM users WHERE username = $1")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var username string
	err = stmt.QueryRow(name).Scan(&username)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}