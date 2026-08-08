# API

Victoria Maler exposes an HTTP JSON API for its crowdfunding and verification features. All endpoints live under `/api`.

## Auth

### `POST /api/auth/register`
Create a new account. Returns a token in the `token` cookie.

```json
{ "name": "Kisumu Cleanup Network", "email": "ngo@example.org", "password": "secret", "role": "ngo" }
```

Roles: `donor` (default), `ngo`, `admin`.

### `POST /api/auth/login`
```json
{ "email": "ngo@example.org", "password": "secret" }
```

### `POST /api/auth/logout`
Clears the session cookie.

## Campaigns

### `POST /api/campaigns` *(auth required)*
```json
{ "title": "Beach Cleanup", "description": "...", "location": "Entebbe", "target_amount": 2500 }
```

### `GET /campaigns.html`
HTML page listing all campaigns.

### `GET /campaign_details.html?id=1`
HTML page with a campaign, its donations, cleanup reports and evidence.

## Donations

### `POST /api/donations`
```json
{ "campaign_id": 1, "donor_name": "Jane", "amount": 50, "message": "Keep going!" }
```
Updates the campaign's `raised_amount` and marks it completed when the target is met.

### `GET /api/campaigns/donations?campaign_id=1`

## Cleanup reports

### `POST /api/cleanups` *(auth required)*
```json
{ "campaign_id": 1, "title": "First clearing", "notes": "...", "location": "Kisumu Pier", "cleanup_date": "2026-07-15", "waste_kg": 850, "volunteers": 40 }
```
The report is hashed and anchored via the Bitcoin anchoring service. The response includes `content_hash`, `txid` and `anchor_status`.

### `GET /api/campaigns/cleanups?campaign_id=1`

### `GET /api/cleanups/verify?id=1`
Recomputes the hash and reports whether the record is intact and anchored.

## Evidence

### `POST /api/evidence?campaign_id=1&report_id=0` *(auth required, multipart/form-data)*
Fields: `file`, `caption`. The file bytes are hashed and anchored.

### `GET /api/campaigns/evidence?campaign_id=1`

## Bitcoin verification

### `POST /api/verify`
```json
{ "hash": "abc123...", "txid": "optional" }
```
Returns `verified`, `method` (`on-chain` or `simulated`) and `note`.

### `GET /api/verify?hash=abc123...`
Looks up the stored anchor for a hash and verifies it.

### `GET /api/verification-records`
Lists every anchored record.

### `GET /api/node-status`
Reports whether the bitcoind node is enabled, reachable, its network, height and wallet balance.

## Dashboard

### `GET /api/dashboard` *(auth required)*
Campaigns, donations and cleanup reports belonging to the signed-in user.

## Response format

All JSON responses follow the same envelope:

```json
{ "success": true, "message": "optional", "data": { } }
```

Errors use the appropriate HTTP status code with `"success": false`.
