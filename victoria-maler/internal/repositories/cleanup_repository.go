package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type CleanupRepository interface {
	Create(r *models.CleanupReport) error
	FindByID(id int64) (*models.CleanupReport, error)
	ListByCampaign(campaignID int64) ([]models.CleanupReport, error)
	ListByUser(userID int64) ([]models.CleanupReport, error)
	UpdateAnchor(id int64, hash, txid, status string) error
}

type SQLiteCleanupRepository struct {
	DB *sql.DB
}

func NewCleanupRepository(db *sql.DB) *SQLiteCleanupRepository {
	return &SQLiteCleanupRepository{DB: db}
}

const cleanupColumns = `id, campaign_id, title, notes, location, cleanup_date, waste_kg, volunteers, content_hash, txid, anchor_status, created_by, created_at`

func (r *SQLiteCleanupRepository) Create(report *models.CleanupReport) error {
	res, err := r.DB.Exec(
		`INSERT INTO cleanup_reports
		 (campaign_id, title, notes, location, cleanup_date, waste_kg, volunteers, content_hash, txid, anchor_status, created_by, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		report.CampaignID, report.Title, report.Notes, report.Location, report.CleanupDate,
		report.WasteKg, report.Volunteers, report.ContentHash, report.TxID, report.AnchorStatus,
		report.CreatedBy, report.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	report.ID = id
	return nil
}

func scanCleanup(row *sql.Row) (*models.CleanupReport, error) {
	var r models.CleanupReport
	if err := row.Scan(&r.ID, &r.CampaignID, &r.Title, &r.Notes, &r.Location, &r.CleanupDate,
		&r.WasteKg, &r.Volunteers, &r.ContentHash, &r.TxID, &r.AnchorStatus, &r.CreatedBy, &r.CreatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}

func scanCleanups(rows *sql.Rows) ([]models.CleanupReport, error) {
	defer rows.Close()
	reports := []models.CleanupReport{}
	for rows.Next() {
		var r models.CleanupReport
		if err := rows.Scan(&r.ID, &r.CampaignID, &r.Title, &r.Notes, &r.Location, &r.CleanupDate,
			&r.WasteKg, &r.Volunteers, &r.ContentHash, &r.TxID, &r.AnchorStatus, &r.CreatedBy, &r.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}

func (r *SQLiteCleanupRepository) FindByID(id int64) (*models.CleanupReport, error) {
	return scanCleanup(r.DB.QueryRow(`SELECT `+cleanupColumns+` FROM cleanup_reports WHERE id = ?`, id))
}

func (r *SQLiteCleanupRepository) ListByCampaign(campaignID int64) ([]models.CleanupReport, error) {
	rows, err := r.DB.Query(`SELECT `+cleanupColumns+` FROM cleanup_reports WHERE campaign_id = ? ORDER BY created_at DESC`, campaignID)
	if err != nil {
		return nil, err
	}
	return scanCleanups(rows)
}

func (r *SQLiteCleanupRepository) ListByUser(userID int64) ([]models.CleanupReport, error) {
	rows, err := r.DB.Query(`SELECT `+cleanupColumns+` FROM cleanup_reports WHERE created_by = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	return scanCleanups(rows)
}

func (r *SQLiteCleanupRepository) UpdateAnchor(id int64, hash, txid, status string) error {
	_, err := r.DB.Exec(
		`UPDATE cleanup_reports SET content_hash = ?, txid = ?, anchor_status = ? WHERE id = ?`,
		hash, txid, status, id,
	)
	return err
}
