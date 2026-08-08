# Database Schema

SQLite database with foreign keys enabled. Migrations are applied in order from
the `migrations/` directory and tracked in `schema_migrations`.

## users

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| name | TEXT | NOT NULL |
| email | TEXT | NOT NULL, UNIQUE |
| password_hash | TEXT | NOT NULL, bcrypt |
| role | TEXT | NOT NULL: `donor` \| `ngo` \| `admin` |
| created_at | INTEGER | unix timestamp |

## campaigns

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| title | TEXT | NOT NULL |
| description | TEXT | default '' |
| location | TEXT | default '' |
| target_amount | REAL | default 0 |
| raised_amount | REAL | default 0 |
| status | TEXT | `active` \| `completed` \| `draft` |
| created_by | INTEGER | FK → users.id |
| created_at | INTEGER | unix timestamp |

## donations

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| campaign_id | INTEGER | FK → campaigns.id |
| user_id | INTEGER | FK → users.id (nullable for anonymous) |
| donor_name | TEXT | default '' |
| amount | REAL | NOT NULL |
| message | TEXT | default '' |
| created_at | INTEGER | unix timestamp |

## cleanup_reports

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| campaign_id | INTEGER | FK → campaigns.id |
| title | TEXT | default '' |
| notes | TEXT | default '' |
| location | TEXT | default '' |
| cleanup_date | TEXT | ISO date |
| waste_kg | REAL | default 0 |
| volunteers | INTEGER | default 0 |
| content_hash | TEXT | SHA-256 hex of the canonical report |
| txid | TEXT | Bitcoin transaction id (or `sim-` prefixed) |
| anchor_status | TEXT | `pending` \| `anchored` \| `simulated` |
| created_by | INTEGER | FK → users.id |
| created_at | INTEGER | unix timestamp |

## evidence

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| report_id | INTEGER | FK → cleanup_reports.id (nullable) |
| campaign_id | INTEGER | FK → campaigns.id |
| file_path | TEXT | NOT NULL |
| file_type | TEXT | MIME type |
| caption | TEXT | default '' |
| content_hash | TEXT | SHA-256 hex of the file bytes |
| txid | TEXT | Bitcoin transaction id |
| anchor_status | TEXT | `pending` \| `anchored` \| `simulated` |
| timestamp | INTEGER | evidence capture time |
| created_at | INTEGER | unix timestamp |

## bitcoin_records

| Column | Type | Notes |
|---|---|---|
| id | INTEGER | PK, autoincrement |
| content_hash | TEXT | NOT NULL, UNIQUE |
| txid | TEXT | default '' |
| op_return | TEXT | hex OP_RETURN script |
| status | TEXT | `anchored` \| `simulated` |
| record_type | TEXT | `cleanup_report` \| `evidence` |
| record_id | INTEGER | id of the referenced record |
| created_at | INTEGER | unix timestamp |

## schema_migrations

Tracks applied migration filenames (`version`) and the time each was applied.
