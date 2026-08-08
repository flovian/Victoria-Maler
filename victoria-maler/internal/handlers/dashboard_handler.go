package handlers

import (
	"net/http"

	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) DashboardPage(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.claims(r)
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	user, err := h.Auth.UserByID(claims.UserID)
	if err != nil {
		utils.InternalError(w, "could not load user")
		return
	}

	myCampaigns, err := h.Campaigns.ListByCreator(claims.UserID)
	if err != nil {
		utils.InternalError(w, "could not load campaigns")
		return
	}

	myDonations, err := h.Donations.ListByUser(claims.UserID)
	if err != nil {
		utils.InternalError(w, "could not load donations")
		return
	}

	myReports, err := h.Cleanups.ListByUser(claims.UserID)
	if err != nil {
		utils.InternalError(w, "could not load cleanup reports")
		return
	}

	totalRaised := 0.0
	for _, d := range myDonations {
		totalRaised += d.Amount
	}

	h.Render.Page(w, r, "dashboard", map[string]interface{}{
		"User":           user,
		"Campaigns":      myCampaigns,
		"Donations":      myDonations,
		"Reports":        myReports,
		"TotalDonated":   totalRaised,
		"CampaignsCount": len(myCampaigns),
		"DonationsCount": len(myDonations),
		"ReportsCount":   len(myReports),
		"EvidencesCount": 0,
	})
}

func (h *HandlerSet) DashboardData(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.claims(r)
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	myCampaigns, _ := h.Campaigns.ListByCreator(claims.UserID)
	myDonations, _ := h.Donations.ListByUser(claims.UserID)
	myReports, _ := h.Cleanups.ListByUser(claims.UserID)

	totalRaised := 0.0
	for _, d := range myDonations {
		totalRaised += d.Amount
	}

	utils.OK(w, map[string]interface{}{
		"campaigns":       myCampaigns,
		"donations":       myDonations,
		"reports":         myReports,
		"total_donated":   totalRaised,
		"campaigns_count": len(myCampaigns),
		"donations_count": len(myDonations),
		"reports_count":   len(myReports),
	})
}
