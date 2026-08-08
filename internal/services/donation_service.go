package services

import (
	"errors"
	"time"

	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
)

var ErrInvalidDonation = errors.New("donation amount must be positive")

type DonationService struct {
	donations repositories.DonationRepository
	campaigns repositories.CampaignRepository
}

func NewDonationService(donations repositories.DonationRepository, campaigns repositories.CampaignRepository) *DonationService {
	return &DonationService{donations: donations, campaigns: campaigns}
}

type DonateInput struct {
	CampaignID int64
	UserID     int64
	DonorName  string
	Amount     float64
	Message    string
}

func (s *DonationService) Donate(in DonateInput) (*models.Donation, error) {
	if in.Amount <= 0 {
		return nil, ErrInvalidDonation
	}

	campaign, err := s.campaigns.FindByID(in.CampaignID)
	if err != nil {
		return nil, ErrCampaignNotFound
	}

	donation := &models.Donation{
		CampaignID: in.CampaignID,
		UserID:     in.UserID,
		DonorName:  in.DonorName,
		Amount:     in.Amount,
		Message:    in.Message,
		CreatedAt:  time.Now().Unix(),
	}
	if err := s.donations.Create(donation); err != nil {
		return nil, err
	}

	newTotal := campaign.RaisedAmount + in.Amount
	if err := s.campaigns.UpdateRaised(in.CampaignID, newTotal); err != nil {
		return nil, err
	}

	if campaign.Status == models.CampaignStatusActive && newTotal >= campaign.TargetAmount {
		if err := s.campaigns.SetStatus(in.CampaignID, models.CampaignStatusCompleted); err != nil {
			return nil, err
		}
	}

	return donation, nil
}

func (s *DonationService) ListByCampaign(campaignID int64) ([]models.Donation, error) {
	return s.donations.ListByCampaign(campaignID)
}

func (s *DonationService) ListByUser(userID int64) ([]models.Donation, error) {
	return s.donations.ListByUser(userID)
}

func (s *DonationService) CampaignTotal(campaignID int64) (float64, error) {
	return s.donations.TotalForCampaign(campaignID)
}
