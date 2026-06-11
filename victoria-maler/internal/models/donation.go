package models

type Donation struct {
    ID         int64   `db:"id" json:"id"`
    CampaignID int64   `db:"campaign_id" json:"campaign_id"`
    UserID     int64   `db:"user_id" json:"user_id"`
    Amount     float64 `db:"amount" json:"amount"`
}
