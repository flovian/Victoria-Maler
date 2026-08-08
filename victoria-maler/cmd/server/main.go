package main

import (
	"log"
	"net/http"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/database"
	"ecochain-victoria/internal/handlers"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/routes"
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
	log.Printf("database ready at %s", cfg.DBPath)

	userRepo := repositories.NewUserRepository(db)
	campaignRepo := repositories.NewCampaignRepository(db)
	donationRepo := repositories.NewDonationRepository(db)
	cleanupRepo := repositories.NewCleanupRepository(db)
	evidenceRepo := repositories.NewEvidenceRepository(db)
	bitcoinRepo := repositories.NewBitcoinRecordRepository(db)

	rpcClient := bitcoin.NewRPCClient(cfg.Bitcoin.RPCURL, cfg.Bitcoin.RPCUser, cfg.Bitcoin.RPCPass)
	bitcoin.SetGlobalClient(rpcClient)

	anchoring := bitcoin.NewAnchoringService(
		cfg.Bitcoin.Enabled,
		rpcClient,
		services.NewBitcoinRecordSaver(bitcoinRepo),
	)

	handlerSet := &handlers.HandlerSet{
		Config:         cfg,
		Render:         handlers.NewRenderer(),
		Auth:           services.NewAuthService(userRepo, cfg),
		Campaigns:      services.NewCampaignService(campaignRepo),
		Donations:      services.NewDonationService(donationRepo, campaignRepo),
		Cleanups:       services.NewCleanupService(cleanupRepo, anchoring),
		Evidences:      services.NewEvidenceService(evidenceRepo, anchoring),
		BitcoinRecords: bitcoinRepo,
		Wallet:         bitcoin.NewWallet(rpcClient),
	}

	server := routes.New(handlerSet, cfg.JWTSecret)

	mux := server.Mux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	addr := ":" + cfg.Port
	log.Printf("Victoria Maler server starting on %s — open http://localhost:%s/", addr, cfg.Port)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
