package models

type CleanupReport struct {
	ID           int64   `db:"id" json:"id"`
	CampaignID   int64   `db:"campaign_id" json:"campaign_id"`
	Title        string  `db:"title" json:"title"`
	Notes        string  `db:"notes" json:"notes"`
	Location     string  `db:"location" json:"location"`
	CleanupDate  string  `db:"cleanup_date" json:"cleanup_date"`
	WasteKg      float64 `db:"waste_kg" json:"waste_kg"`
	Volunteers   int     `db:"volunteers" json:"volunteers"`
	ContentHash  string  `db:"content_hash" json:"content_hash"`
	TxID         string  `db:"txid" json:"txid"`
	AnchorStatus string  `db:"anchor_status" json:"anchor_status"`
	CreatedBy    int64   `db:"created_by" json:"created_by"`
	CreatedAt    int64   `db:"created_at" json:"created_at"`
}

const (
	AnchorStatusAnchored  = "anchored"
	AnchorStatusPending   = "pending"
	AnchorStatusSimulated = "simulated"
)
