package postgres

import (
	"database/sql"
	"errors"
	"user-service/models"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(user *models.User) (*models.User, error) {
	id := uuid.New()
	query := `INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id, name, email, password, created_at, updated_at`
	err := r.db.QueryRow(query, id, user.Name, user.Email, user.Password).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (r *PostgresRepository) GetUser(id string) (*models.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresRepository) UpdateUser(user *models.User) (*models.User, error) {
	query := `UPDATE users SET name = $2, email = $3, password = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING id, name, email, password, created_at, updated_at`
	err := r.db.QueryRow(query, user.ID, user.Name, user.Email, user.Password).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (r *PostgresRepository) DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *PostgresRepository) ListUsers(limit, page int) ([]*models.User, error) {
	if page == 1 {
		page = 0
	}
	query := `SELECT id, name, email, password, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, page*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
