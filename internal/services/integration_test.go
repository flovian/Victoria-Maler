package services_test

import (
	"database/sql"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/database"
	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/services"
)

func newTestStack(t *testing.T) (auth *services.AuthService, campaigns *services.CampaignService, donations *services.DonationService, cleanups *services.CleanupService, evidences *services.EvidenceService, db *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db, "../../migrations"); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{JWTSecret: "integration-test-secret"}

	userRepo := repositories.NewUserRepository(db)
	campaignRepo := repositories.NewCampaignRepository(db)
	donationRepo := repositories.NewDonationRepository(db)
	cleanupRepo := repositories.NewCleanupRepository(db)
	evidenceRepo := repositories.NewEvidenceRepository(db)
	bitcoinRepo := repositories.NewBitcoinRecordRepository(db)

	anchoring := bitcoin.NewAnchoringService(false, nil, services.NewBitcoinRecordSaver(bitcoinRepo))

	return services.NewAuthService(userRepo, cfg),
		services.NewCampaignService(campaignRepo),
		services.NewDonationService(donationRepo, campaignRepo),
		services.NewCleanupService(cleanupRepo, anchoring),
		services.NewEvidenceService(evidenceRepo, anchoring),
		db
}

func TestFullPlatformFlow(t *testing.T) {
	auth, campaigns, donations, cleanups, evidences, db := newTestStack(t)

	ngo, err := auth.Register("Cleanup NGO", "ngo@test.org", "password123", models.RoleNGO)
	if err != nil {
		t.Fatalf("register ngo: %v", err)
	}
	if ngo.ID == 0 {
		t.Fatal("expected generated user id")
	}

	token, err := auth.Token(ngo)
	if err != nil || token == "" {
		t.Fatalf("token generation failed: %v", err)
	}
	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if claims.UserID != ngo.ID {
		t.Fatalf("token user mismatch: %d != %d", claims.UserID, ngo.ID)
	}

	campaign, err := campaigns.CreateCampaign(services.CreateCampaignInput{
		Title:        "Entebbe Beach Cleanup",
		Description:  "Remove shoreline plastic waste.",
		Location:     "Entebbe, Uganda",
		TargetAmount: 2500,
		CreatedBy:    ngo.ID,
	})
	if err != nil {
		t.Fatalf("create campaign: %v", err)
	}

	donor, err := auth.Register("Donor Jane", "donor@test.org", "password123", models.RoleDonor)
	if err != nil {
		t.Fatalf("register donor: %v", err)
	}

	donation, err := donations.Donate(services.DonateInput{
		CampaignID: campaign.ID,
		UserID:     donor.ID,
		DonorName:  "Donor Jane",
		Amount:     100,
	})
	if err != nil {
		t.Fatalf("donate: %v", err)
	}
	if donation.Amount != 100 {
		t.Fatalf("donation amount mismatch")
	}

	updated, err := campaigns.GetCampaign(campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if updated.RaisedAmount != 100 {
		t.Fatalf("expected raised 100, got %.2f", updated.RaisedAmount)
	}

	report, err := cleanups.SubmitReport(services.SubmitCleanupInput{
		CampaignID:  campaign.ID,
		Title:       "First cleanup",
		Notes:       "Collected plastic and hyacinth.",
		Location:    "Entebbe Pier",
		CleanupDate: "2026-07-10",
		WasteKg:     420,
		Volunteers:  25,
		CreatedBy:   ngo.ID,
	})
	if err != nil {
		t.Fatalf("submit cleanup: %v", err)
	}
	if len(report.ContentHash) != 64 {
		t.Fatalf("expected 64-char content hash, got %d", len(report.ContentHash))
	}
	if !strings.HasPrefix(report.TxID, "sim-") {
		t.Fatalf("expected simulated txid without a node, got %q", report.TxID)
	}
	if report.AnchorStatus != models.AnchorStatusSimulated {
		t.Fatalf("expected simulated anchor status, got %q", report.AnchorStatus)
	}

	verification := cleanups.VerifyReport(report)
	if !verification.Verified {
		t.Fatalf("intact report should verify: %+v", verification)
	}

	_, err = evidences.SubmitEvidence(services.SubmitEvidenceInput{
		CampaignID: campaign.ID,
		ReportID:   report.ID,
		FilePath:   "before.jpg",
		FileType:   "image/jpeg",
		Caption:    "Shoreline before cleanup",
		FileData:   []byte("fake-jpeg-bytes"),
		Timestamp:  report.CreatedAt,
	})
	if err != nil {
		t.Fatalf("submit evidence: %v", err)
	}

	evidenceList, err := evidences.ListByCampaign(campaign.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(evidenceList) != 1 {
		t.Fatalf("expected 1 evidence item, got %d", len(evidenceList))
	}

	// Tamper detection: alter the waste figure directly in the database.
	if _, err := db.Exec(`UPDATE cleanup_reports SET waste_kg = 9999 WHERE id = ?`, report.ID); err != nil {
		t.Fatal(err)
	}
	tampered, err := cleanups.GetReport(report.ID)
	if err != nil {
		t.Fatal(err)
	}
	tamperedVerification := cleanups.VerifyReport(tampered)
	if tamperedVerification.Verified {
		t.Fatalf("tampered report must fail verification: %+v", tamperedVerification)
	}
}
