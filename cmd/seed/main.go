package main

import (
	"log"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/database"
	"ecochain-victoria/internal/models"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db, "migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	userRepo := repositories.NewUserRepository(db)
	campaignRepo := repositories.NewCampaignRepository(db)
	donationRepo := repositories.NewDonationRepository(db)
	cleanupRepo := repositories.NewCleanupRepository(db)
	bitcoinRepo := repositories.NewBitcoinRecordRepository(db)

	auth := services.NewAuthService(userRepo, cfg)
	campaigns := services.NewCampaignService(campaignRepo)
	donations := services.NewDonationService(donationRepo, campaignRepo)

	rpcClient := bitcoin.NewRPCClient(cfg.Bitcoin.RPCURL, cfg.Bitcoin.RPCUser, cfg.Bitcoin.RPCPass)
	anchoring := bitcoin.NewAnchoringService(cfg.Bitcoin.Enabled, rpcClient, services.NewBitcoinRecordSaver(bitcoinRepo))
	cleanups := services.NewCleanupService(cleanupRepo, anchoring)

	if _, err := auth.Register("Platform Admin", "admin@victoriamaler.org", "admin12345", models.RoleAdmin); err != nil {
		log.Printf("admin user: %v", err)
	}
	ngo, err := auth.Register("Kisumu Cleanup Network", "ngo@victoriamaler.org", "ngo12345", models.RoleNGO)
	if err != nil {
		log.Printf("ngo user: %v", err)
	}

	if _, err := campaigns.CreateCampaign(services.CreateCampaignInput{
		Title:        "Beach Cleanup - Entebbe",
		Description:  "Remove plastic waste from Entebbe's shoreline with local volunteers and youth groups.",
		Location:     "Entebbe, Uganda",
		TargetAmount: 2500,
		CreatedBy:    ngo.ID,
	}); err != nil {
		log.Printf("campaign 1: %v", err)
	}

	beach, err := campaigns.CreateCampaign(services.CreateCampaignInput{
		Title:        "Water Hyacinth Removal - Kisumu",
		Description:  "Clear water hyacinth mats blocking the lakeshore and restore fishing access.",
		Location:     "Kisumu, Kenya",
		TargetAmount: 5000,
		CreatedBy:    ngo.ID,
	})
	if err != nil {
		log.Printf("campaign 2: %v", err)
	}

	if _, err := donations.Donate(services.DonateInput{
		CampaignID: beach.ID,
		UserID:     1,
		DonorName:  "Anonymous Donor",
		Amount:     120.50,
		Message:    "Keep up the great work protecting Lake Victoria!",
	}); err != nil {
		log.Printf("donation: %v", err)
	}

	if _, err := cleanups.SubmitReport(services.SubmitCleanupInput{
		CampaignID:  beach.ID,
		Title:       "First hyacinth clearing",
		Notes:       "Volunteers cleared hyacinth from the main jetty area.",
		Location:    "Kisumu Pier",
		CleanupDate: "2026-07-15",
		WasteKg:     850,
		Volunteers:  40,
		CreatedBy:   ngo.ID,
	}); err != nil {
		log.Printf("cleanup report: %v", err)
	}

	log.Println("demo data seeded successfully")
	log.Println("  admin: admin@victoriamaler.org / admin12345")
	log.Println("  ngo:   ngo@victoriamaler.org / ngo12345")
}
