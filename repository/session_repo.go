package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
)

type SessionRepository struct {
	db *db.DB
}

func NewSessionRepository(db *db.DB) *SessionRepository {

	return &SessionRepository{db: db}

}

func (r *SessionRepository) Create(sessionID string, accountID int, expiresAt time.Time) (*models.Session, error) {

	query := `
	INSERT INTO sessions (id,account_id,expires_at)
	VALUES ($1,$2,$3)
	RETURNING id, account_id, expires_at, created_at
	`

	session := &models.Session{}

	err := r.db.QueryRow(query, sessionID, accountID, expiresAt).Scan(&session.ID, session.AccountID, session.CreatedAt, session.ExpiresAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create a session: %w", err)
	}
	return session, nil
}

func (r *SessionRepository) GetByID(sessionID string) (*models.Session, error) {
	query := `
	SELECT id,account_id,expires_at,created_at
	FROM sessions
	WHERE id = $1
	`
	session := &models.Session{}
	err := r.db.QueryRow(query, sessionID).Scan(&session.ID, session.AccountID, session.CreatedAt, session.ExpiresAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (r *SessionRepository) GetByAccountID(accountID string) ([]*models.Session, error) {
	sessions := make([]*models.Session, 0)
	query := `
	SELECT id, account_id, expires_at, created_at
	FROM sessions
	WHERE account_id = $1
	`
	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		session := &models.Session{}

		err := rows.Scan(session.AccountID, session.ID, session.CreatedAt, session.ExpiresAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) Delete(sessionID string) error {
	query := `
	DELETE FROM sessions WHERE id = $1
	`
	result, err := r.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *SessionRepository) DeleteAccountByID(accountID int) error {

	query := `
	DELETE FROM sessions WHERE account_id = $1
	`

	rows, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *SessionRepository) DeleteExpired() (int, error) {
	query := `
	DELETE FROM sessions WHERE expires_at < $1
	`
	result, err := r.db.Exec(query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check delete result: %w", err)
	}

	return int(rowsAffected), nil

}

func (r *SessionRepository) IsValid(sessionID string) (bool, error) {
	session, err := r.GetByID(sessionID)
	if err != nil {
		return false, nil
	}

	if session.IsExpired() {
		r.Delete(sessionID)
		return false, nil
	}

	return true, nil
}
