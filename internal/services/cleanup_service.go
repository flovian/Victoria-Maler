package services

import (
	"errors"
	"fmt"
	"time"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/utils"
)

var ErrCleanupInput = errors.New("cleanup report requires a campaign and location")

type CleanupService struct {
	cleanups  repositories.CleanupRepository
	anchoring *bitcoin.AnchoringService
}

func NewCleanupService(cleanups repositories.CleanupRepository, anchoring *bitcoin.AnchoringService) *CleanupService {
	return &CleanupService{cleanups: cleanups, anchoring: anchoring}
}

type SubmitCleanupInput struct {
	CampaignID  int64
	Title       string
	Notes       string
	Location    string
	CleanupDate string
	WasteKg     float64
	Volunteers  int
	CreatedBy   int64
}

func (s *CleanupService) SubmitReport(in SubmitCleanupInput) (*models.CleanupReport, error) {
	if in.CampaignID <= 0 || in.Location == "" {
		return nil, ErrCleanupInput
	}

	hash := utils.HashHex(canonicalCleanup(in))

	report := &models.CleanupReport{
		CampaignID:   in.CampaignID,
		Title:        in.Title,
		Notes:        in.Notes,
		Location:     in.Location,
		CleanupDate:  in.CleanupDate,
		WasteKg:      in.WasteKg,
		Volunteers:   in.Volunteers,
		ContentHash:  hash,
		AnchorStatus: models.AnchorStatusPending,
		CreatedBy:    in.CreatedBy,
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.cleanups.Create(report); err != nil {
		return nil, err
	}

	anchor, err := s.anchoring.AnchorHash(hash, models.RecordTypeCleanup, report.ID)
	if err != nil {
		return nil, err
	}
	if err := s.cleanups.UpdateAnchor(report.ID, hash, anchor.TxID, anchor.Status); err != nil {
		return nil, err
	}
	report.TxID = anchor.TxID
	report.AnchorStatus = anchor.Status
	return report, nil
}

func (s *CleanupService) ListByCampaign(campaignID int64) ([]models.CleanupReport, error) {
	return s.cleanups.ListByCampaign(campaignID)
}

func (s *CleanupService) ListByUser(userID int64) ([]models.CleanupReport, error) {
	return s.cleanups.ListByUser(userID)
}

func (s *CleanupService) GetReport(id int64) (*models.CleanupReport, error) {
	return s.cleanups.FindByID(id)
}

// VerifyReport recomputes the canonical hash and checks whether the stored
// hash still matches, i.e. whether the report has been tampered with.
func (s *CleanupService) VerifyReport(report *models.CleanupReport) *bitcoin.VerificationResult {
	recomputed := utils.HashHex(canonicalCleanup(SubmitCleanupInput{
		CampaignID:  report.CampaignID,
		Title:       report.Title,
		Notes:       report.Notes,
		Location:    report.Location,
		CleanupDate: report.CleanupDate,
		WasteKg:     report.WasteKg,
		Volunteers:  report.Volunteers,
	}))
	if recomputed != report.ContentHash {
		return &bitcoin.VerificationResult{
			Verified: false,
			TxID:     report.TxID,
			Hash:     report.ContentHash,
			Note:     "Record has been modified since submission - hash no longer matches.",
		}
	}
	return bitcoin.VerifyHash(report.ContentHash, report.TxID)
}

func canonicalCleanup(in SubmitCleanupInput) string {
	return fmt.Sprintf("cleanup|%d|%s|%s|%s|%s|%.4f|%d",
		in.CampaignID, in.Title, in.Location, in.CleanupDate, in.Notes, in.WasteKg, in.Volunteers)
}
