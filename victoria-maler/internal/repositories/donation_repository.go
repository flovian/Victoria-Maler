package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type DonationRepository interface {
	Create(d *models.Donation) error
	ListByCampaign(campaignID int64) ([]models.Donation, error)
	ListByUser(userID int64) ([]models.Donation, error)
	TotalForCampaign(campaignID int64) (float64, error)
}

type SQLiteDonationRepository struct {
	DB *sql.DB
}

func NewDonationRepository(db *sql.DB) *SQLiteDonationRepository {
	return &SQLiteDonationRepository{DB: db}
}

const donationColumns = `id, campaign_id, user_id, donor_name, amount, message, created_at`

func (r *SQLiteDonationRepository) Create(d *models.Donation) error {
	var userID interface{}
	if d.UserID > 0 {
		userID = d.UserID
	}
	res, err := r.DB.Exec(
		`INSERT INTO donations (campaign_id, user_id, donor_name, amount, message, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		d.CampaignID, userID, d.DonorName, d.Amount, d.Message, d.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

func scanDonation(row *sql.Row) (*models.Donation, error) {
	var d models.Donation
	if err := row.Scan(&d.ID, &d.CampaignID, &d.UserID, &d.DonorName, &d.Amount, &d.Message, &d.CreatedAt); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *SQLiteDonationRepository) ListByCampaign(campaignID int64) ([]models.Donation, error) {
	rows, err := r.DB.Query(`SELECT `+donationColumns+` FROM donations WHERE campaign_id = ? ORDER BY created_at DESC`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	donations := []models.Donation{}
	for rows.Next() {
		var d models.Donation
		if err := rows.Scan(&d.ID, &d.CampaignID, &d.UserID, &d.DonorName, &d.Amount, &d.Message, &d.CreatedAt); err != nil {
			return nil, err
		}
		donations = append(donations, d)
	}
	return donations, rows.Err()
}

func (r *SQLiteDonationRepository) ListByUser(userID int64) ([]models.Donation, error) {
	rows, err := r.DB.Query(`SELECT `+donationColumns+` FROM donations WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	donations := []models.Donation{}
	for rows.Next() {
		var d models.Donation
		if err := rows.Scan(&d.ID, &d.CampaignID, &d.UserID, &d.DonorName, &d.Amount, &d.Message, &d.CreatedAt); err != nil {
			return nil, err
		}
		donations = append(donations, d)
	}
	return donations, rows.Err()
}

func (r *SQLiteDonationRepository) TotalForCampaign(campaignID int64) (float64, error) {
	var total float64
	err := r.DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM donations WHERE campaign_id = ?`, campaignID).Scan(&total)
	return total, err
}
