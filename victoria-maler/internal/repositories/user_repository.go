package repositories

import "ecochain-victoria/internal/models"

type UserRepository interface {
    Create(u *models.User) error
}
