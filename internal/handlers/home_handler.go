package handlers

import (
	"net/http"

	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) Home(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.Campaigns.ListCampaigns()
	if err != nil {
		utils.InternalError(w, "could not load campaigns")
		return
	}

	activeCount := 0
	for _, c := range campaigns {
		if c.Status == "active" {
			activeCount++
		}
	}

	totalRaised := 0.0
	for _, c := range campaigns {
		totalRaised += c.RaisedAmount
	}

	h.Render.Page(w, r, "home", map[string]interface{}{
		"Campaigns":     campaigns,
		"ActiveCount":   activeCount,
		"CampaignCount": len(campaigns),
		"TotalRaised":   totalRaised,
	})
}

func (h *HandlerSet) About(w http.ResponseWriter, r *http.Request) {
	h.Render.Page(w, r, "about", nil)
}
