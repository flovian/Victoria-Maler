package models

type Evidence struct {
	ID           int64  `db:"id" json:"id"`
	ReportID     int64  `db:"report_id" json:"report_id"`
	CampaignID   int64  `db:"campaign_id" json:"campaign_id"`
	FilePath     string `db:"file_path" json:"file_path"`
	FileType     string `db:"file_type" json:"file_type"`
	Caption      string `db:"caption" json:"caption"`
	ContentHash  string `db:"content_hash" json:"content_hash"`
	TxID         string `db:"txid" json:"txid"`
	AnchorStatus string `db:"anchor_status" json:"anchor_status"`
	Timestamp    int64  `db:"timestamp" json:"timestamp"`
	CreatedAt    int64  `db:"created_at" json:"created_at"`
}
