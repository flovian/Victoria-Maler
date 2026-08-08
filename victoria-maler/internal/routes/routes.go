package routes

import (
	"net/http"

	"ecochain-victoria/internal/handlers"
	"ecochain-victoria/internal/middleware"
)

type Server struct {
	handler *http.ServeMux
	h       *handlers.HandlerSet
	secret  string
}

func New(h *handlers.HandlerSet, secret string) *Server {
	s := &Server{
		handler: http.NewServeMux(),
		h:       h,
		secret:  secret,
	}
	s.register()
	return s
}

func (s *Server) Handler() http.Handler {
	return middleware.Logger(s.handler)
}

// Mux exposes the underlying mux so the server can mount static assets.
func (s *Server) Mux() *http.ServeMux {
	return s.handler
}

func (s *Server) register() {
	h := s.h
	mux := s.handler
	optAuth := middleware.OptAuth(s.secret)

	// Pages
	mux.Handle("/", optAuth(http.HandlerFunc(h.Home)))
	mux.Handle("GET /home", optAuth(http.HandlerFunc(h.Home)))
	mux.Handle("GET /about.html", optAuth(http.HandlerFunc(h.About)))
	mux.Handle("GET /campaigns.html", optAuth(http.HandlerFunc(h.CampaignsPage)))
	mux.Handle("GET /campaign_details.html", optAuth(http.HandlerFunc(h.CampaignDetailsPage)))
	mux.Handle("GET /donate.html", optAuth(http.HandlerFunc(h.DonatePage)))
	mux.Handle("GET /evidence_tracker.html", optAuth(http.HandlerFunc(h.EvidenceTrackerPage)))
	mux.Handle("GET /bitcoin_verify.html", optAuth(http.HandlerFunc(h.BitcoinVerifyPage)))
	mux.Handle("GET /login.html", optAuth(http.HandlerFunc(h.LoginPage)))
	mux.Handle("GET /register.html", optAuth(http.HandlerFunc(h.RegisterPage)))
	mux.Handle("GET /dashboard.html", middleware.RequireAuth(s.secret)(http.HandlerFunc(h.DashboardPage)))
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/logout", h.Logout)

	// Campaigns
	mux.Handle("POST /api/campaigns", middleware.RequireAuth(s.secret)(http.HandlerFunc(h.CreateCampaign)))

	// Donations
	mux.HandleFunc("POST /api/donations", h.Donate)
	mux.HandleFunc("GET /api/campaigns/donations", h.CampaignDonations)

	// Cleanup reports
	mux.Handle("POST /api/cleanups", middleware.RequireAuth(s.secret)(http.HandlerFunc(h.SubmitCleanup)))
	mux.HandleFunc("GET /api/campaigns/cleanups", h.CampaignCleanups)
	mux.HandleFunc("GET /api/cleanups/verify", h.VerifyCleanup)

	// Evidence
	mux.Handle("POST /api/evidence", middleware.RequireAuth(s.secret)(http.HandlerFunc(h.SubmitEvidence)))
	mux.HandleFunc("GET /api/campaigns/evidence", h.CampaignEvidence)

	// Bitcoin verification
	mux.HandleFunc("POST /api/verify", h.VerifyHash)
	mux.HandleFunc("GET /api/verify", h.VerifyByHashOnly)
	mux.HandleFunc("GET /api/verification-records", h.VerificationRecords)
	mux.HandleFunc("GET /api/node-status", h.NodeStatus)

	// Dashboard API
	mux.Handle("GET /api/dashboard", middleware.RequireAuth(s.secret)(http.HandlerFunc(h.DashboardData)))
}
