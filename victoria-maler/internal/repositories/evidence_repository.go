package repositories

import "ecochain-victoria/internal/models"

type EvidenceRepository interface {
    Save(e *models.Evidence) error
}
