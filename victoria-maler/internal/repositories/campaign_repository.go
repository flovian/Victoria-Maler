package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type CampaignRepository interface {
	Create(c *models.Campaign) error
	List() ([]models.Campaign, error)
	FindByID(id int64) (*models.Campaign, error)
	ListByCreator(userID int64) ([]models.Campaign, error)
	UpdateRaised(id int64, amount float64) error
	SetStatus(id int64, status string) error
}

type SQLiteCampaignRepository struct {
	DB *sql.DB
}

func NewCampaignRepository(db *sql.DB) *SQLiteCampaignRepository {
	return &SQLiteCampaignRepository{DB: db}
}

const campaignColumns = `id, title, description, location, target_amount, raised_amount, status, created_by, created_at`

func (r *SQLiteCampaignRepository) Create(c *models.Campaign) error {
	res, err := r.DB.Exec(
		`INSERT INTO campaigns (title, description, location, target_amount, raised_amount, status, created_by, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		c.Title, c.Description, c.Location, c.TargetAmount, c.RaisedAmount, c.Status, c.CreatedBy, c.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}

func scanCampaign(row *sql.Row) (*models.Campaign, error) {
	var c models.Campaign
	if err := row.Scan(&c.ID, &c.Title, &c.Description, &c.Location, &c.TargetAmount,
		&c.RaisedAmount, &c.Status, &c.CreatedBy, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SQLiteCampaignRepository) List() ([]models.Campaign, error) {
	rows, err := r.DB.Query(`SELECT ` + campaignColumns + ` FROM campaigns ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := []models.Campaign{}
	for rows.Next() {
		var c models.Campaign
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Location, &c.TargetAmount,
			&c.RaisedAmount, &c.Status, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, rows.Err()
}

func (r *SQLiteCampaignRepository) FindByID(id int64) (*models.Campaign, error) {
	return scanCampaign(r.DB.QueryRow(`SELECT `+campaignColumns+` FROM campaigns WHERE id = ?`, id))
}

func (r *SQLiteCampaignRepository) ListByCreator(userID int64) ([]models.Campaign, error) {
	rows, err := r.DB.Query(`SELECT `+campaignColumns+` FROM campaigns WHERE created_by = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := []models.Campaign{}
	for rows.Next() {
		var c models.Campaign
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Location, &c.TargetAmount,
			&c.RaisedAmount, &c.Status, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, rows.Err()
}

func (r *SQLiteCampaignRepository) UpdateRaised(id int64, amount float64) error {
	_, err := r.DB.Exec(
		`UPDATE campaigns SET raised_amount = ? WHERE id = ?`,
		amount, id,
	)
	return err
}

func (r *SQLiteCampaignRepository) SetStatus(id int64, status string) error {
	_, err := r.DB.Exec(
		`UPDATE campaigns SET status = ? WHERE id = ?`,
		status, id,
	)
	return err
}
