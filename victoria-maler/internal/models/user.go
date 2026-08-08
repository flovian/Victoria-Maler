package models

type User struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Role         string `db:"role" json:"role"`
	CreatedAt    int64  `db:"created_at" json:"created_at"`
}

const (
	RoleDonor = "donor"
	RoleNGO   = "ngo"
	RoleAdmin = "admin"
)
