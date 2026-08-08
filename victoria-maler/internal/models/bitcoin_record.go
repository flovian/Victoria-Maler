package models

type BitcoinRecord struct {
	ID          int64  `db:"id" json:"id"`
	ContentHash string `db:"content_hash" json:"content_hash"`
	TxID        string `db:"txid" json:"txid"`
	OpReturn    string `db:"op_return" json:"op_return"`
	Status      string `db:"status" json:"status"`
	RecordType  string `db:"record_type" json:"record_type"`
	RecordID    int64  `db:"record_id" json:"record_id"`
	CreatedAt   int64  `db:"created_at" json:"created_at"`
}

const (
	RecordTypeCleanup  = "cleanup_report"
	RecordTypeEvidence = "evidence"
)
