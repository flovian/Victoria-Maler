package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type UserRepository interface {
	Create(u *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id int64) (*models.User, error)
}

type SQLiteUserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{DB: db}
}

func (r *SQLiteUserRepository) Create(u *models.User) error {
	res, err := r.DB.Exec(
		`INSERT INTO users (name, email, password_hash, role, created_at) VALUES (?, ?, ?, ?, ?)`,
		u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *SQLiteUserRepository) FindByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE email = ?`,
		email,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *SQLiteUserRepository) FindByID(id int64) (*models.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, name, email, password_hash, role, created_at FROM users WHERE id = ?`,
		id,
	)
	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}
