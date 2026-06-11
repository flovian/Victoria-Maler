package repositories

import "ecochain-victoria/internal/models"

type CleanupRepository interface {
    Create(r *models.CleanupReport) error
}
