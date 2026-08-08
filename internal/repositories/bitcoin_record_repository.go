package repositories

import (
	"database/sql"

	"ecochain-victoria/internal/models"
)

type BitcoinRecordRepository interface {
	Create(r *models.BitcoinRecord) error
	FindByHash(contentHash string) (*models.BitcoinRecord, error)
	List() ([]models.BitcoinRecord, error)
	UpdateStatus(id int64, txid, status string) error
}

type SQLiteBitcoinRecordRepository struct {
	DB *sql.DB
}

func NewBitcoinRecordRepository(db *sql.DB) *SQLiteBitcoinRecordRepository {
	return &SQLiteBitcoinRecordRepository{DB: db}
}

const bitcoinColumns = `id, content_hash, txid, op_return, status, record_type, record_id, created_at`

func (r *SQLiteBitcoinRecordRepository) Create(record *models.BitcoinRecord) error {
	res, err := r.DB.Exec(
		`INSERT INTO bitcoin_records (content_hash, txid, op_return, status, record_type, record_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		record.ContentHash, record.TxID, record.OpReturn, record.Status, record.RecordType, record.RecordID, record.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	record.ID = id
	return nil
}

func scanBitcoinRecord(row *sql.Row) (*models.BitcoinRecord, error) {
	var r models.BitcoinRecord
	if err := row.Scan(&r.ID, &r.ContentHash, &r.TxID, &r.OpReturn, &r.Status, &r.RecordType, &r.RecordID, &r.CreatedAt); err != nil {
		return nil, err
	}
	return &r, nil
}

func (r *SQLiteBitcoinRecordRepository) FindByHash(contentHash string) (*models.BitcoinRecord, error) {
	return scanBitcoinRecord(r.DB.QueryRow(`SELECT `+bitcoinColumns+` FROM bitcoin_records WHERE content_hash = ?`, contentHash))
}

func (r *SQLiteBitcoinRecordRepository) List() ([]models.BitcoinRecord, error) {
	rows, err := r.DB.Query(`SELECT ` + bitcoinColumns + ` FROM bitcoin_records ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []models.BitcoinRecord{}
	for rows.Next() {
		var rec models.BitcoinRecord
		if err := rows.Scan(&rec.ID, &rec.ContentHash, &rec.TxID, &rec.OpReturn, &rec.Status, &rec.RecordType, &rec.RecordID, &rec.CreatedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

func (r *SQLiteBitcoinRecordRepository) UpdateStatus(id int64, txid, status string) error {
	_, err := r.DB.Exec(
		`UPDATE bitcoin_records SET txid = ?, status = ? WHERE id = ?`,
		txid, status, id,
	)
	return err
}
