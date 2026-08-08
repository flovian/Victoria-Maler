-- 006_bitcoin_records.sql
CREATE TABLE IF NOT EXISTS bitcoin_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  content_hash TEXT NOT NULL UNIQUE,
  txid TEXT NOT NULL DEFAULT '',
  op_return TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'pending',
  record_type TEXT NOT NULL DEFAULT '',
  record_id INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL
);
