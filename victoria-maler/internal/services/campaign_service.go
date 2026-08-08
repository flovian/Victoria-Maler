package services

import (
	"errors"
	"time"

	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
)

var ErrCampaignNotFound = errors.New("campaign not found")

type CampaignService struct {
	campaigns repositories.CampaignRepository
}

func NewCampaignService(campaigns repositories.CampaignRepository) *CampaignService {
	return &CampaignService{campaigns: campaigns}
}

type CreateCampaignInput struct {
	Title        string
	Description  string
	Location     string
	TargetAmount float64
	CreatedBy    int64
}

func (s *CampaignService) CreateCampaign(in CreateCampaignInput) (*models.Campaign, error) {
	if in.Title == "" {
		return nil, errors.New("campaign title is required")
	}
	if in.TargetAmount <= 0 {
		return nil, errors.New("campaign target amount must be positive")
	}

	campaign := &models.Campaign{
		Title:        in.Title,
		Description:  in.Description,
		Location:     in.Location,
		TargetAmount: in.TargetAmount,
		Status:       models.CampaignStatusActive,
		CreatedBy:    in.CreatedBy,
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.campaigns.Create(campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func (s *CampaignService) ListCampaigns() ([]models.Campaign, error) {
	return s.campaigns.List()
}

func (s *CampaignService) GetCampaign(id int64) (*models.Campaign, error) {
	return s.campaigns.FindByID(id)
}

func (s *CampaignService) ListByCreator(userID int64) ([]models.Campaign, error) {
	return s.campaigns.ListByCreator(userID)
}
