-- 005_evidence.sql
CREATE TABLE IF NOT EXISTS evidence (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  report_id INTEGER,
  campaign_id INTEGER NOT NULL,
  file_path TEXT NOT NULL,
  file_type TEXT NOT NULL DEFAULT '',
  caption TEXT NOT NULL DEFAULT '',
  content_hash TEXT NOT NULL DEFAULT '',
  txid TEXT NOT NULL DEFAULT '',
  anchor_status TEXT NOT NULL DEFAULT 'pending',
  timestamp INTEGER NOT NULL,
  created_at INTEGER NOT NULL,
  FOREIGN KEY (report_id) REFERENCES cleanup_reports(id),
  FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);
