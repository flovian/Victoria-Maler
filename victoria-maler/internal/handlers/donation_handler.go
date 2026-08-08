package handlers

import (
	"net/http"

	"ecochain-victoria/internal/services"
	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) DonatePage(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.Campaigns.ListCampaigns()
	if err != nil {
		utils.InternalError(w, "could not load campaigns")
		return
	}
	h.Render.Page(w, r, "donate", map[string]interface{}{
		"Campaigns": campaigns,
	})
}

type donateRequest struct {
	CampaignID int64   `json:"campaign_id"`
	DonorName  string  `json:"donor_name"`
	Amount     float64 `json:"amount"`
	Message    string  `json:"message"`
}

func (h *HandlerSet) Donate(w http.ResponseWriter, r *http.Request) {
	var req donateRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}

	in := services.DonateInput{
		CampaignID: req.CampaignID,
		DonorName:  req.DonorName,
		Amount:     req.Amount,
		Message:    req.Message,
	}
	if claims, ok := h.claims(r); ok {
		in.UserID = claims.UserID
	}

	donation, err := h.Donations.Donate(in)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, donation)
}

func (h *HandlerSet) CampaignDonations(w http.ResponseWriter, r *http.Request) {
	campaignID, err := queryID(r, "campaign_id")
	if err != nil {
		utils.BadRequest(w, "campaign_id required")
		return
	}
	donations, err := h.Donations.ListByCampaign(campaignID)
	if err != nil {
		utils.InternalError(w, "could not load donations")
		return
	}
	utils.OK(w, donations)
}
