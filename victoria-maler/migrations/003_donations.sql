-- 003_donations.sql
CREATE TABLE IF NOT EXISTS donations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  campaign_id INTEGER,
  user_id INTEGER,
  amount REAL
);
