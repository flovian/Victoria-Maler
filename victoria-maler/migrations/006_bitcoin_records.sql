-- 006_bitcoin_records.sql
CREATE TABLE IF NOT EXISTS bitcoin_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  txid TEXT,
  op_return TEXT,
  created_at INTEGER
);
