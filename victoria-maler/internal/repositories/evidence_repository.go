package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type EvidenceRepository interface {
	Save(e *models.Evidence) error
	FindByID(id int64) (*models.Evidence, error)
	ListByCampaign(campaignID int64) ([]models.Evidence, error)
	ListByReport(reportID int64) ([]models.Evidence, error)
	UpdateAnchor(id int64, txid, status string) error
}

type SQLiteEvidenceRepository struct {
	DB *sql.DB
}

func NewEvidenceRepository(db *sql.DB) *SQLiteEvidenceRepository {
	return &SQLiteEvidenceRepository{DB: db}
}

const evidenceColumns = `id, report_id, campaign_id, file_path, file_type, caption, content_hash, txid, anchor_status, timestamp, created_at`

func (r *SQLiteEvidenceRepository) Save(e *models.Evidence) error {
	var reportID interface{}
	if e.ReportID > 0 {
		reportID = e.ReportID
	}
	res, err := r.DB.Exec(
		`INSERT INTO evidence
		 (report_id, campaign_id, file_path, file_type, caption, content_hash, txid, anchor_status, timestamp, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		reportID, e.CampaignID, e.FilePath, e.FileType, e.Caption, e.ContentHash,
		e.TxID, e.AnchorStatus, e.Timestamp, e.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func scanEvidence(row *sql.Row) (*models.Evidence, error) {
	var e models.Evidence
	if err := row.Scan(&e.ID, &e.ReportID, &e.CampaignID, &e.FilePath, &e.FileType, &e.Caption,
		&e.ContentHash, &e.TxID, &e.AnchorStatus, &e.Timestamp, &e.CreatedAt); err != nil {
		return nil, err
	}
	return &e, nil
}

func scanEvidences(rows *sql.Rows) ([]models.Evidence, error) {
	defer rows.Close()
	evidences := []models.Evidence{}
	for rows.Next() {
		var e models.Evidence
		if err := rows.Scan(&e.ID, &e.ReportID, &e.CampaignID, &e.FilePath, &e.FileType, &e.Caption,
			&e.ContentHash, &e.TxID, &e.AnchorStatus, &e.Timestamp, &e.CreatedAt); err != nil {
			return nil, err
		}
		evidences = append(evidences, e)
	}
	return evidences, rows.Err()
}

func (r *SQLiteEvidenceRepository) FindByID(id int64) (*models.Evidence, error) {
	return scanEvidence(r.DB.QueryRow(`SELECT `+evidenceColumns+` FROM evidence WHERE id = ?`, id))
}

func (r *SQLiteEvidenceRepository) ListByCampaign(campaignID int64) ([]models.Evidence, error) {
	rows, err := r.DB.Query(`SELECT `+evidenceColumns+` FROM evidence WHERE campaign_id = ? ORDER BY created_at DESC`, campaignID)
	if err != nil {
		return nil, err
	}
	return scanEvidences(rows)
}

func (r *SQLiteEvidenceRepository) ListByReport(reportID int64) ([]models.Evidence, error) {
	rows, err := r.DB.Query(`SELECT `+evidenceColumns+` FROM evidence WHERE report_id = ? ORDER BY created_at DESC`, reportID)
	if err != nil {
		return nil, err
	}
	return scanEvidences(rows)
}

func (r *SQLiteEvidenceRepository) UpdateAnchor(id int64, txid, status string) error {
	_, err := r.DB.Exec(
		`UPDATE evidence SET txid = ?, anchor_status = ? WHERE id = ?`,
		txid, status, id,
	)
	return err
}
