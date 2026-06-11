package repositories

import "ecochain-victoria/internal/models"

type CampaignRepository interface {
    List() ([]models.Campaign, error)
}
