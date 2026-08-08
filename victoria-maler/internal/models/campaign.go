package models

type Campaign struct {
	ID           int64   `db:"id" json:"id"`
	Title        string  `db:"title" json:"title"`
	Description  string  `db:"description" json:"description"`
	Location     string  `db:"location" json:"location"`
	TargetAmount float64 `db:"target_amount" json:"target_amount"`
	RaisedAmount float64 `db:"raised_amount" json:"raised_amount"`
	Status       string  `db:"status" json:"status"`
	CreatedBy    int64   `db:"created_by" json:"created_by"`
	CreatedAt    int64   `db:"created_at" json:"created_at"`
}

const (
	CampaignStatusActive    = "active"
	CampaignStatusCompleted = "completed"
	CampaignStatusDraft     = "draft"
)
