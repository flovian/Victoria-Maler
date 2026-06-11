package models

type BitcoinRecord struct {
    ID         int64  `db:"id" json:"id"`
    TxID       string `db:"txid" json:"txid"`
    OpReturn   string `db:"op_return" json:"op_return"`
    CreatedAt  int64  `db:"created_at" json:"created_at"`
}
