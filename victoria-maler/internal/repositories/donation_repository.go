package repositories

import "ecochain-victoria/internal/models"

type DonationRepository interface {
    Create(d *models.Donation) error
}
