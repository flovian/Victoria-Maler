package models

type CleanupReport struct {
    ID        int64  `db:"id" json:"id"`
    CampaignID int64 `db:"campaign_id" json:"campaign_id"`
    Notes     string `db:"notes" json:"notes"`
}
