-- 004_cleanup_reports.sql
CREATE TABLE IF NOT EXISTS cleanup_reports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  campaign_id INTEGER NOT NULL,
  title TEXT NOT NULL DEFAULT '',
  notes TEXT NOT NULL DEFAULT '',
  location TEXT NOT NULL DEFAULT '',
  cleanup_date TEXT NOT NULL DEFAULT '',
  waste_kg REAL NOT NULL DEFAULT 0,
  volunteers INTEGER NOT NULL DEFAULT 0,
  content_hash TEXT NOT NULL DEFAULT '',
  txid TEXT NOT NULL DEFAULT '',
  anchor_status TEXT NOT NULL DEFAULT 'pending',
  created_by INTEGER,
  created_at INTEGER NOT NULL,
  FOREIGN KEY (campaign_id) REFERENCES campaigns(id),
  FOREIGN KEY (created_by) REFERENCES users(id)
);
