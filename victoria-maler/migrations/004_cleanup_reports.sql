-- 004_cleanup_reports.sql
CREATE TABLE IF NOT EXISTS cleanup_reports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  campaign_id INTEGER,
  notes TEXT
);
