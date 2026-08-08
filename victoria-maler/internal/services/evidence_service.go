package services

import (
	"errors"
	"time"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/utils"
)

var ErrEvidenceInput = errors.New("evidence requires a campaign and file")

type EvidenceService struct {
	evidences repositories.EvidenceRepository
	anchoring *bitcoin.AnchoringService
}

func NewEvidenceService(evidences repositories.EvidenceRepository, anchoring *bitcoin.AnchoringService) *EvidenceService {
	return &EvidenceService{evidences: evidences, anchoring: anchoring}
}

type SubmitEvidenceInput struct {
	CampaignID int64
	ReportID   int64
	FilePath   string
	FileType   string
	Caption    string
	FileData   []byte
	Timestamp  int64
}

func (s *EvidenceService) SubmitEvidence(in SubmitEvidenceInput) (*models.Evidence, error) {
	if in.CampaignID <= 0 || in.FilePath == "" {
		return nil, ErrEvidenceInput
	}

	hash := utils.HashHexBytes(in.FileData)

	evidence := &models.Evidence{
		ReportID:     in.ReportID,
		CampaignID:   in.CampaignID,
		FilePath:     in.FilePath,
		FileType:     in.FileType,
		Caption:      in.Caption,
		ContentHash:  hash,
		AnchorStatus: models.AnchorStatusPending,
		Timestamp:    in.Timestamp,
		CreatedAt:    time.Now().Unix(),
	}
	if err := s.evidences.Save(evidence); err != nil {
		return nil, err
	}

	anchor, err := s.anchoring.AnchorHash(hash, models.RecordTypeEvidence, evidence.ID)
	if err != nil {
		return nil, err
	}
	if err := s.evidences.UpdateAnchor(evidence.ID, anchor.TxID, anchor.Status); err != nil {
		return nil, err
	}
	evidence.TxID = anchor.TxID
	evidence.AnchorStatus = anchor.Status
	return evidence, nil
}

func (s *EvidenceService) ListByCampaign(campaignID int64) ([]models.Evidence, error) {
	return s.evidences.ListByCampaign(campaignID)
}

func (s *EvidenceService) ListByReport(reportID int64) ([]models.Evidence, error) {
	return s.evidences.ListByReport(reportID)
}

func (s *EvidenceService) GetEvidence(id int64) (*models.Evidence, error) {
	return s.evidences.FindByID(id)
}

// VerifyEvidence checks whether the file bytes hash to the anchored hash.
func (s *EvidenceService) VerifyEvidence(evidence *models.Evidence, fileData []byte) *bitcoin.VerificationResult {
	recomputed := utils.HashHexBytes(fileData)
	if recomputed != evidence.ContentHash {
		return &bitcoin.VerificationResult{
			Verified: false,
			TxID:     evidence.TxID,
			Hash:     evidence.ContentHash,
			Note:     "File has been altered since submission - hash no longer matches.",
		}
	}
	return bitcoin.VerifyHash(evidence.ContentHash, evidence.TxID)
}

// bitcoinRecordSaver adapts the bitcoin anchoring service to persist
// anchor records through the repository layer.
type bitcoinRecordSaver struct {
	records repositories.BitcoinRecordRepository
}

func NewBitcoinRecordSaver(records repositories.BitcoinRecordRepository) bitcoin.RecordSaver {
	return &bitcoinRecordSaver{records: records}
}

func (s *bitcoinRecordSaver) Save(record *bitcoin.AnchorRecord) error {
	return s.records.Create(&models.BitcoinRecord{
		ContentHash: record.ContentHash,
		TxID:        record.TxID,
		OpReturn:    record.OpReturn,
		Status:      record.Status,
		RecordType:  record.RecordType,
		RecordID:    record.RecordID,
		CreatedAt:   record.CreatedAt,
	})
}
