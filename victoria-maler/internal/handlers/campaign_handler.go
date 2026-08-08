package handlers

import (
	"net/http"
	"strconv"

	"ecochain-victoria/internal/services"
	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) CampaignsPage(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.Campaigns.ListCampaigns()
	if err != nil {
		utils.InternalError(w, "could not load campaigns")
		return
	}
	h.Render.Page(w, r, "campaigns", map[string]interface{}{
		"Campaigns": campaigns,
	})
}

func (h *HandlerSet) CampaignDetailsPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		utils.NotFound(w, "campaign not found")
		return
	}
	campaign, err := h.Campaigns.GetCampaign(id)
	if err != nil {
		utils.NotFound(w, "campaign not found")
		return
	}
	donations, err := h.Donations.ListByCampaign(id)
	if err != nil {
		utils.InternalError(w, "could not load donations")
		return
	}
	reports, err := h.Cleanups.ListByCampaign(id)
	if err != nil {
		utils.InternalError(w, "could not load cleanup reports")
		return
	}
	evidences, err := h.Evidences.ListByCampaign(id)
	if err != nil {
		utils.InternalError(w, "could not load evidence")
		return
	}
	h.Render.Page(w, r, "campaign_details", map[string]interface{}{
		"Campaign":  campaign,
		"Donations": donations,
		"Reports":   reports,
		"Evidences": evidences,
	})
}

type createCampaignRequest struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Location     string  `json:"location"`
	TargetAmount float64 `json:"target_amount"`
}

func (h *HandlerSet) CreateCampaign(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.claims(r)
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var req createCampaignRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}

	campaign, err := h.Campaigns.CreateCampaign(services.CreateCampaignInput{
		Title:        req.Title,
		Description:  req.Description,
		Location:     req.Location,
		TargetAmount: req.TargetAmount,
		CreatedBy:    claims.UserID,
	})
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, campaign)
}
